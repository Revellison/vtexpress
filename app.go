package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"vtgui/internal"
	"vtgui/internal/models"
	"vtgui/internal/modules"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	orchestrator *internal.Orchestrator
	storage      *modules.Storage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	storagePath, err := resolveStoragePath()
	if err != nil {
		runtime.LogErrorf(ctx, "resolve storage path: %v", err)
		return
	}

	storage, err := modules.NewStorage(storagePath)
	if err != nil {
		runtime.LogErrorf(ctx, "init storage: %v", err)
		return
	}

	a.storage = storage
	a.orchestrator = internal.NewOrchestrator(storage)
}

func (a *App) shutdown(_ context.Context) {
	if a.storage != nil {
		_ = a.storage.Close()
	}
}

func (a *App) SaveSettings(request models.SaveSettingsRequest) (models.SaveSettingsResponse, error) {
	if a.orchestrator == nil {
		return models.SaveSettingsResponse{}, errors.New("backend is not initialized")
	}
	return a.orchestrator.SaveSettings(a.ctx, request)
}

func (a *App) GetSettings() (models.SettingsState, error) {
	if a.orchestrator == nil {
		return models.SettingsState{}, errors.New("backend is not initialized")
	}
	return a.orchestrator.GetSettings(a.ctx)
}

func (a *App) AnalyzeFileByPath(request models.AnalyzeByPathRequest) (models.ScanResult, error) {
	if a.orchestrator == nil {
		return models.ScanResult{}, errors.New("backend is not initialized")
	}
	return a.orchestrator.AnalyzeFileByPath(a.ctx, request)
}

func (a *App) AnalyzeFileByContent(request models.AnalyzeByContentRequest) (models.ScanResult, error) {
	if a.orchestrator == nil {
		return models.ScanResult{}, errors.New("backend is not initialized")
	}
	return a.orchestrator.AnalyzeFileByContent(a.ctx, request)
}

func (a *App) SummarizePayload(request models.AISummaryRequest) (string, error) {
	if a.orchestrator == nil {
		return "", errors.New("backend is not initialized")
	}
	return a.orchestrator.SummarizePayload(a.ctx, request)
}

func (a *App) GetHistory() ([]models.HistoryItem, error) {
	if a.orchestrator == nil {
		return nil, errors.New("backend is not initialized")
	}
	return a.orchestrator.ListHistory(a.ctx)
}

func (a *App) PickFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}
	return path, nil
}

func resolveStoragePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}

	dir := filepath.Join(configDir, "vtexpress")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("create storage dir: %w", err)
	}

	return filepath.Join(dir, "vtexpress.db"), nil
}
