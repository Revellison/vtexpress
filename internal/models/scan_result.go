package models

import "encoding/json"

type AIProvider string

const (
	AIProviderGemini     AIProvider = "gemini"
	AIProviderOpenRouter AIProvider = "openrouter"
)

type SettingsState struct {
	HasVTAPIKey   bool       `json:"hasVtApiKey"`
	HasAIAPIKey   bool       `json:"hasAiApiKey"`
	AIProvider    AIProvider `json:"aiProvider"`
	AIModel       string     `json:"aiModel"`
	AutoAISummary bool       `json:"autoAiSummary"`
}

type SaveSettingsRequest struct {
	VTAPIKey      string     `json:"vtApiKey"`
	AIAPIKey      string     `json:"aiApiKey"`
	AIProvider    AIProvider `json:"aiProvider"`
	AIModel       string     `json:"aiModel"`
	AutoAISummary bool       `json:"autoAiSummary"`
}

type SaveSettingsResponse struct {
	VTValid bool   `json:"vtValid"`
	AIValid bool   `json:"aiValid"`
	Message string `json:"message"`
}

type VTAnalysisEngine struct {
	EngineName    string `json:"engineName"`
	Category      string `json:"category"`
	Result        string `json:"result,omitempty"`
	Method        string `json:"method,omitempty"`
	EngineUpdate  string `json:"engineUpdate,omitempty"`
	EngineVersion string `json:"engineVersion,omitempty"`
}

type EnginesVerdict struct {
	Engines []VTAnalysisEngine `json:"engines"`
}

type IPInfo struct {
	Address  string `json:"address"`
	ASN      string `json:"asn,omitempty"`
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
}

type DomainInfo struct {
	Domain   string `json:"domain"`
	Resolved string `json:"resolved,omitempty"`
	ASN      string `json:"asn,omitempty"`
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
}

type IPDomainReport struct {
	IPs     []IPInfo     `json:"ip"`
	Domains []DomainInfo `json:"domain"`
}

type VTSummary struct {
	Verdict         string `json:"verdict"`
	TotalEngines    int    `json:"totalEngines"`
	Malicious       int    `json:"malicious"`
	Suspicious      int    `json:"suspicious"`
	Undetected      int    `json:"undetected"`
	Harmless        int    `json:"harmless"`
	Timeout         int    `json:"timeout"`
	Failure         int    `json:"failure"`
	TypeUnsupported int    `json:"typeUnsupported"`
	FileName        string `json:"fileName,omitempty"`
	FileType        string `json:"fileType,omitempty"`
	SHA256          string `json:"sha256,omitempty"`
	Size            int64  `json:"size,omitempty"`
}

type ScanPayload struct {
	EnginesVerdict EnginesVerdict `json:"engines_verdict"`
	IPDomain       IPDomainReport `json:"ip_domain"`
	VTSummary      VTSummary      `json:"vt_summary"`
}

type ScanResult struct {
	FileName   string          `json:"fileName"`
	FileSHA256 string          `json:"fileSha256"`
	ScannedAt  string          `json:"scannedAt"`
	Payload    ScanPayload     `json:"payload"`
	RawVT      json.RawMessage `json:"rawVt"`
	AISummary  string          `json:"aiSummary,omitempty"`
}

type AnalyzeByPathRequest struct {
	Path         string `json:"path"`
	RunAISummary bool   `json:"runAiSummary"`
}

type AnalyzeByContentRequest struct {
	FileName     string `json:"fileName"`
	Base64Data   string `json:"base64Data"`
	RunAISummary bool   `json:"runAiSummary"`
}

type AISummaryRequest struct {
	Payload ScanPayload `json:"payload"`
}
