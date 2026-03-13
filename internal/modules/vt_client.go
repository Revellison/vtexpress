package modules

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	vtBaseURL              = "https://www.virustotal.com/api/v3"
	vtDefaultPollInterval  = 5 * time.Second
	vtDefaultPollTimeout   = 4 * time.Minute
	vtMaxRequestsPerMinute = 4
)

type VTClient struct {
	httpClient   *http.Client
	limiter      <-chan time.Time
	limiterSetup sync.Once
}

func NewVTClient() *VTClient {
	return &VTClient{
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

func (c *VTClient) initLimiter() {
	c.limiterSetup.Do(func() {
		interval := time.Minute / vtMaxRequestsPerMinute
		ticker := time.NewTicker(interval)
		c.limiter = ticker.C
	})
}

func (c *VTClient) waitRateLimit(ctx context.Context) error {
	c.initLimiter()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.limiter:
		return nil
	}
}

func (c *VTClient) ValidateAPIKey(ctx context.Context, apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return errors.New("vt api key is empty")
	}
	emptyHash := sha256.Sum256(nil)
	_, status, _, err := c.request(ctx, http.MethodGet, "/files/"+hex.EncodeToString(emptyHash[:]), apiKey, nil, "")
	if err != nil {
		return err
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return errors.New("vt api key is invalid")
	}
	if status >= 200 && status < 500 {
		return nil
	}
	return fmt.Errorf("vt api key validation failed: status %d", status)
}

func (c *VTClient) ComputeSHA256(fileContent []byte) string {
	sum := sha256.Sum256(fileContent)
	return hex.EncodeToString(sum[:])
}

func (c *VTClient) GetFileReport(ctx context.Context, apiKey, sha256Hex string) ([]byte, int, error) {
	body, status, _, err := c.request(ctx, http.MethodGet, "/files/"+sha256Hex, apiKey, nil, "")
	return body, status, err
}

func (c *VTClient) UploadFileForAnalysis(ctx context.Context, apiKey, fileName string, fileContent []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("create multipart file: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		return "", fmt.Errorf("write multipart file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	bodyRaw, status, _, err := c.request(ctx, http.MethodPost, "/files", apiKey, &body, writer.FormDataContentType())
	if err != nil {
		return "", err
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("upload file failed: status %d", status)
	}

	var response struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyRaw, &response); err != nil {
		return "", fmt.Errorf("decode upload response: %w", err)
	}
	if response.Data.ID == "" {
		return "", errors.New("upload response did not return analysis id")
	}

	return response.Data.ID, nil
}

func (c *VTClient) WaitForAnalysisAndFetchReport(ctx context.Context, apiKey, analysisID string) ([]byte, error) {
	pollCtx, cancel := context.WithTimeout(ctx, vtDefaultPollTimeout)
	defer cancel()

	ticker := time.NewTicker(vtDefaultPollInterval)
	defer ticker.Stop()

	for {
		analysisRaw, status, _, err := c.request(pollCtx, http.MethodGet, "/analyses/"+analysisID, apiKey, nil, "")
		if err != nil {
			return nil, err
		}
		if status < 200 || status >= 300 {
			return nil, fmt.Errorf("analysis poll failed: status %d", status)
		}

		var analysis struct {
			Data struct {
				Attributes struct {
					Status string `json:"status"`
				} `json:"attributes"`
			} `json:"data"`
			Meta struct {
				FileInfo struct {
					SHA256 string `json:"sha256"`
				} `json:"file_info"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(analysisRaw, &analysis); err != nil {
			return nil, fmt.Errorf("decode analysis: %w", err)
		}

		if strings.EqualFold(analysis.Data.Attributes.Status, "completed") && analysis.Meta.FileInfo.SHA256 != "" {
			reportRaw, reportStatus, err := c.GetFileReport(pollCtx, apiKey, analysis.Meta.FileInfo.SHA256)
			if err != nil {
				return nil, err
			}
			if reportStatus < 200 || reportStatus >= 300 {
				return nil, fmt.Errorf("fetch final report failed: status %d", reportStatus)
			}
			return reportRaw, nil
		}

		select {
		case <-pollCtx.Done():
			return nil, pollCtx.Err()
		case <-ticker.C:
		}
	}
}

func (c *VTClient) request(
	ctx context.Context,
	method string,
	path string,
	apiKey string,
	body io.Reader,
	contentType string,
) ([]byte, int, http.Header, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, 0, nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, vtBaseURL+path, body)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("create vt request: %w", err)
	}

	request.Header.Set("x-apikey", apiKey)
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("execute vt request: %w", err)
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, response.Header, fmt.Errorf("read vt response: %w", err)
	}

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return raw, response.StatusCode, response.Header, errors.New("vt authorization failed")
	}

	return raw, response.StatusCode, response.Header, nil
}
