package modules

import (
	"context"
	"fmt"
	"sort"

	"vtgui/internal/models"
)

type EnginesReportModule struct{}

func NewEnginesReportModule() *EnginesReportModule {
	return &EnginesReportModule{}
}

func (m *EnginesReportModule) Build(_ context.Context, vtData map[string]any) (models.EnginesVerdict, error) {
	attributes, err := nestedMap(vtData, "data", "attributes")
	if err != nil {
		return models.EnginesVerdict{}, err
	}

	resultsRaw, ok := attributes["last_analysis_results"]
	if !ok {
		return models.EnginesVerdict{Engines: []models.VTAnalysisEngine{}}, nil
	}

	results, ok := resultsRaw.(map[string]any)
	if !ok {
		return models.EnginesVerdict{}, fmt.Errorf("invalid last_analysis_results format")
	}

	engines := make([]models.VTAnalysisEngine, 0, len(results))
	for engineKey, entryRaw := range results {
		entry, ok := entryRaw.(map[string]any)
		if !ok {
			continue
		}

		engine := models.VTAnalysisEngine{
			EngineName:    pickString(entry, "engine_name", engineKey),
			Category:      pickString(entry, "category", "unknown"),
			Result:        pickString(entry, "result", ""),
			Method:        pickString(entry, "method", ""),
			EngineUpdate:  pickString(entry, "engine_update", ""),
			EngineVersion: pickString(entry, "engine_version", ""),
		}
		engines = append(engines, engine)
	}

	sort.Slice(engines, func(i, j int) bool {
		return engines[i].EngineName < engines[j].EngineName
	})

	return models.EnginesVerdict{Engines: engines}, nil
}
