package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"vtgui/internal/models"
)

const (
	geminiBaseURL    = "https://generativelanguage.googleapis.com/v1beta"
	openRouterAPIURL = "https://openrouter.ai/api/v1"
)

type AIClient struct {
	httpClient *http.Client
}

func NewAIClient() *AIClient {
	return &AIClient{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *AIClient) ValidateAPIKey(ctx context.Context, provider models.AIProvider, apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return errors.New("ai api key is empty")
	}

	switch provider {
	case models.AIProviderGemini:
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			geminiBaseURL+"/models?key="+apiKey,
			nil,
		)
		if err != nil {
			return err
		}
		response, err := c.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
			return errors.New("gemini api key is invalid")
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return fmt.Errorf("gemini validation failed: status %d", response.StatusCode)
		}
		return nil
	case models.AIProviderOpenRouter:
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, openRouterAPIURL+"/models", nil)
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+apiKey)
		response, err := c.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
			return errors.New("openrouter api key is invalid")
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return fmt.Errorf("openrouter validation failed: status %d", response.StatusCode)
		}
		return nil
	default:
		return fmt.Errorf("unsupported ai provider: %s", provider)
	}
}

func (c *AIClient) Summarize(
	ctx context.Context,
	provider models.AIProvider,
	model string,
	apiKey string,
	payload models.ScanPayload,
) (string, error) {
	if strings.TrimSpace(model) == "" {
		return "", errors.New("ai model is empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return "", errors.New("ai api key is empty")
	}

	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal payload for ai: %w", err)
	}

	systemPrompt := buildAISystemPrompt()
	userPrompt := "Сформируй структурированное резюме по payload VirusTotal ниже. Включи: 1) итоговый вердикт, 2) важные срабатывания движков, 3) сетевую часть ip/domain с рисками, 4) краткие рекомендации пользователю.\\n\\nPayload:\\n" + string(payloadJSON)

	switch provider {
	case models.AIProviderGemini:
		return c.summarizeGemini(ctx, model, apiKey, systemPrompt, userPrompt)
	case models.AIProviderOpenRouter:
		return c.summarizeOpenRouter(ctx, model, apiKey, systemPrompt, userPrompt)
	default:
		return "", fmt.Errorf("unsupported ai provider: %s", provider)
	}
}

func (c *AIClient) summarizeGemini(ctx context.Context, model, apiKey, systemPrompt, userPrompt string) (string, error) {
	body := map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{{"text": systemPrompt}},
		},
		"contents": []map[string]any{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": userPrompt}},
			},
		},
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", geminiBaseURL, model, apiKey)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rawBody))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("gemini summarize failed: status %d body %s", response.StatusCode, string(bodyBytes))
	}

	var payload struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}

	if len(payload.Candidates) == 0 || len(payload.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("gemini returned empty summary")
	}

	return strings.TrimSpace(payload.Candidates[0].Content.Parts[0].Text), nil
}

func (c *AIClient) summarizeOpenRouter(ctx context.Context, model, apiKey, systemPrompt, userPrompt string) (string, error) {
	body := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterAPIURL+"/chat/completions", bytes.NewReader(rawBody))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("openrouter summarize failed: status %d body %s", response.StatusCode, string(bodyBytes))
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}

	if len(payload.Choices) == 0 {
		return "", errors.New("openrouter returned empty summary")
	}

	return strings.TrimSpace(payload.Choices[0].Message.Content), nil
}

func buildAISystemPrompt() string {
	return "Ты ассистент по анализу безопасности файлов. Работай строго на основе переданного payload VirusTotal. Формат ответа: \n" +
		"1) Итоговый вердикт (кратко, 1-2 предложения).\n" +
		"2) Детали по движкам: сколько вредоносных/подозрительных/чистых, перечисли только важные срабатывания.\n" +
		"3) Сетевая активность: IP и домены, укажи ASN/провайдера/регион и потенциальные риски.\n" +
		"4) Практические рекомендации пользователю (3-6 пунктов).\n" +
		"Если данных недостаточно, явно напиши что информации недостаточно, не выдумывай факты. Пиши на русском языке."
}
