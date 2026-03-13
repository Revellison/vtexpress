package modules

import (
	"context"
	"vtgui/internal/models"
)

type VTSummaryModule struct{}

func NewVTSummaryModule() *VTSummaryModule {
	return &VTSummaryModule{}
}

func (m *VTSummaryModule) Build(_ context.Context, vtData map[string]any) (models.VTSummary, error) {
	attributes, err := nestedMap(vtData, "data", "attributes")
	if err != nil {
		return models.VTSummary{}, err
	}

	stats, _ := attributes["last_analysis_stats"].(map[string]any)
	malicious := pickInt(stats, "malicious")
	suspicious := pickInt(stats, "suspicious")
	undetected := pickInt(stats, "undetected")
	harmless := pickInt(stats, "harmless")
	timeout := pickInt(stats, "timeout") + pickInt(stats, "confirmed-timeout")
	failure := pickInt(stats, "failure")
	typeUnsupported := pickInt(stats, "type-unsupported")

	totalEngines := malicious + suspicious + undetected + harmless + timeout + failure + typeUnsupported
	verdict := "clean"
	if malicious > 0 || suspicious > 0 {
		verdict = "malicious"
	}

	fileName := ""
	names, _ := attributes["names"].([]any)
	if len(names) > 0 {
		if first, ok := names[0].(string); ok {
			fileName = first
		}
	}

	return models.VTSummary{
		Verdict:         verdict,
		TotalEngines:    totalEngines,
		Malicious:       malicious,
		Suspicious:      suspicious,
		Undetected:      undetected,
		Harmless:        harmless,
		Timeout:         timeout,
		Failure:         failure,
		TypeUnsupported: typeUnsupported,
		FileName:        fileName,
		FileType:        pickString(attributes, "type_description", ""),
		SHA256:          pickString(attributes, "sha256", ""),
		Size:            pickInt64(attributes, "size"),
	}, nil
}
