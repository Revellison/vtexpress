package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vtgui/internal/models"
	"vtgui/internal/modules"
)

const (
	settingAIProvider    = "ai_provider"
	settingAIModel       = "ai_model"
	settingAutoAISummary = "auto_ai_summary"
)

type Orchestrator struct {
	vtClient *modules.VTClient
	reporter *modules.ReportBuilder
	aiClient *modules.AIClient
	storage  *modules.Storage
	keyring  *modules.KeyringStore
}

func NewOrchestrator(storage *modules.Storage) *Orchestrator {
	return &Orchestrator{
		vtClient: NewSafeVTClient(),
		reporter: modules.NewReportBuilder(),
		aiClient: modules.NewAIClient(),
		storage:  storage,
		keyring:  modules.NewKeyringStore(),
	}
}

func NewSafeVTClient() *modules.VTClient {
	return modules.NewVTClient()
}

func (o *Orchestrator) GetSettings(ctx context.Context) (models.SettingsState, error) {
	vtKey, err := o.keyring.GetVTAPIKey()
	if err != nil {
		return models.SettingsState{}, err
	}
	aiKey, err := o.keyring.GetAIAPIKey()
	if err != nil {
		return models.SettingsState{}, err
	}

	provider, err := o.storage.GetSetting(settingAIProvider)
	if err != nil {
		return models.SettingsState{}, err
	}
	if provider == "" {
		provider = string(models.AIProviderGemini)
	}

	model, err := o.storage.GetSetting(settingAIModel)
	if err != nil {
		return models.SettingsState{}, err
	}
	autoValue, err := o.storage.GetSetting(settingAutoAISummary)
	if err != nil {
		return models.SettingsState{}, err
	}

	return models.SettingsState{
		HasVTAPIKey:   vtKey != "",
		HasAIAPIKey:   aiKey != "",
		AIProvider:    models.AIProvider(provider),
		AIModel:       model,
		AutoAISummary: strings.EqualFold(autoValue, "true"),
	}, nil
}

func (o *Orchestrator) SaveSettings(ctx context.Context, request models.SaveSettingsRequest) (models.SaveSettingsResponse, error) {
	provider := request.AIProvider
	if provider == "" {
		provider = models.AIProviderGemini
	}

	if strings.TrimSpace(request.AIModel) == "" {
		return models.SaveSettingsResponse{}, errors.New("ai model is required")
	}

	vtValid := false
	aiValid := false

	if err := o.vtClient.ValidateAPIKey(ctx, request.VTAPIKey); err != nil {
		return models.SaveSettingsResponse{VTValid: false, AIValid: false, Message: err.Error()}, nil
	}
	vtValid = true

	if err := o.aiClient.ValidateAPIKey(ctx, provider, request.AIAPIKey); err != nil {
		return models.SaveSettingsResponse{VTValid: true, AIValid: false, Message: err.Error()}, nil
	}
	aiValid = true

	if err := o.keyring.SetVTAPIKey(strings.TrimSpace(request.VTAPIKey)); err != nil {
		return models.SaveSettingsResponse{}, err
	}
	if err := o.keyring.SetAIAPIKey(strings.TrimSpace(request.AIAPIKey)); err != nil {
		return models.SaveSettingsResponse{}, err
	}

	if err := o.storage.SetSetting(settingAIProvider, string(provider), "string"); err != nil {
		return models.SaveSettingsResponse{}, err
	}
	if err := o.storage.SetSetting(settingAIModel, strings.TrimSpace(request.AIModel), "string"); err != nil {
		return models.SaveSettingsResponse{}, err
	}
	if err := o.storage.SetSetting(settingAutoAISummary, fmt.Sprintf("%t", request.AutoAISummary), "bool"); err != nil {
		return models.SaveSettingsResponse{}, err
	}

	return models.SaveSettingsResponse{
		VTValid: vtValid,
		AIValid: aiValid,
		Message: "API keys validated and securely saved",
	}, nil
}

func (o *Orchestrator) AnalyzeFileByPath(ctx context.Context, request models.AnalyzeByPathRequest) (models.ScanResult, error) {
	if strings.TrimSpace(request.Path) == "" {
		return models.ScanResult{}, errors.New("file path is empty")
	}

	fileData, err := os.ReadFile(request.Path)
	if err != nil {
		return models.ScanResult{}, fmt.Errorf("read file: %w", err)
	}
	return o.processFile(ctx, filepath.Base(request.Path), fileData, request.RunAISummary)
}

func (o *Orchestrator) AnalyzeFileByContent(ctx context.Context, request models.AnalyzeByContentRequest) (models.ScanResult, error) {
	if strings.TrimSpace(request.FileName) == "" {
		return models.ScanResult{}, errors.New("file name is empty")
	}

	raw, err := base64.StdEncoding.DecodeString(request.Base64Data)
	if err != nil {
		return models.ScanResult{}, fmt.Errorf("decode base64 file: %w", err)
	}

	return o.processFile(ctx, request.FileName, raw, request.RunAISummary)
}

func (o *Orchestrator) SummarizePayload(ctx context.Context, request models.AISummaryRequest) (string, error) {
	state, err := o.GetSettings(ctx)
	if err != nil {
		return "", err
	}
	aiKey, err := o.keyring.GetAIAPIKey()
	if err != nil {
		return "", err
	}
	if aiKey == "" {
		return "", errors.New("ai api key is not configured")
	}

	return o.aiClient.Summarize(ctx, state.AIProvider, state.AIModel, aiKey, request.Payload)
}

func (o *Orchestrator) ListHistory(_ context.Context) ([]models.HistoryItem, error) {
	return o.storage.ListScans(300)
}

func (o *Orchestrator) processFile(ctx context.Context, fileName string, fileData []byte, forceAISummary bool) (models.ScanResult, error) {
	vtKey, err := o.keyring.GetVTAPIKey()
	if err != nil {
		return models.ScanResult{}, err
	}
	if vtKey == "" {
		return models.ScanResult{}, errors.New("virustotal api key is not configured")
	}

	sha := o.vtClient.ComputeSHA256(fileData)
	reportRaw, status, err := o.vtClient.GetFileReport(ctx, vtKey, sha)
	if err != nil && status != 404 {
		return models.ScanResult{}, err
	}

	if status == 404 {
		analysisID, uploadErr := o.vtClient.UploadFileForAnalysis(ctx, vtKey, fileName, fileData)
		if uploadErr != nil {
			return models.ScanResult{}, uploadErr
		}
		reportRaw, err = o.vtClient.WaitForAnalysisAndFetchReport(ctx, vtKey, analysisID)
		if err != nil {
			return models.ScanResult{}, err
		}
	}

	payload, err := o.reporter.BuildPayload(ctx, reportRaw)
	if err != nil {
		return models.ScanResult{}, err
	}

	scannedAt := time.Now().UTC().Format(time.RFC3339)
	result := models.ScanResult{
		FileName:   fileName,
		FileSHA256: sha,
		ScannedAt:  scannedAt,
		Payload:    payload,
		RawVT:      json.RawMessage(reportRaw),
	}

	state, err := o.GetSettings(ctx)
	if err != nil {
		return models.ScanResult{}, err
	}
	shouldSummarize := state.AutoAISummary || forceAISummary
	if shouldSummarize {
		aiKey, err := o.keyring.GetAIAPIKey()
		if err != nil {
			return models.ScanResult{}, err
		}
		if aiKey != "" {
			summary, err := o.aiClient.Summarize(ctx, state.AIProvider, state.AIModel, aiKey, payload)
			if err == nil {
				result.AISummary = summary
			}
		}
	}

	if _, err := o.storage.SaveScan(result); err != nil {
		return models.ScanResult{}, err
	}

	return result, nil
}
