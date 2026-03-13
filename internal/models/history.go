package models

import "encoding/json"

type HistoryItem struct {
	ID         int64           `json:"id"`
	CreatedAt  string          `json:"createdAt"`
	FileName   string          `json:"fileName"`
	FileSHA256 string          `json:"fileSha256"`
	Payload    ScanPayload     `json:"payload"`
	RawVT      json.RawMessage `json:"rawVt"`
	AISummary  string          `json:"aiSummary,omitempty"`
}
