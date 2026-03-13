package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"vtgui/internal/models"
)

type ReportBuilder struct {
	enginesModule  *EnginesReportModule
	ipDomainModule *IPDomainReportModule
	summaryModule  *VTSummaryModule
}

func NewReportBuilder() *ReportBuilder {
	return &ReportBuilder{
		enginesModule:  NewEnginesReportModule(),
		ipDomainModule: NewIPDomainReportModule(),
		summaryModule:  NewVTSummaryModule(),
	}
}

func (b *ReportBuilder) BuildPayload(ctx context.Context, vtRaw []byte) (models.ScanPayload, error) {
	var vtData map[string]any
	if err := json.Unmarshal(vtRaw, &vtData); err != nil {
		return models.ScanPayload{}, fmt.Errorf("decode vt report: %w", err)
	}

	var (
		payload  models.ScanPayload
		wg       sync.WaitGroup
		errMu    sync.Mutex
		buildErr error
	)

	setErr := func(err error) {
		errMu.Lock()
		defer errMu.Unlock()
		if buildErr == nil {
			buildErr = err
		}
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		engines, err := b.enginesModule.Build(ctx, vtData)
		if err != nil {
			setErr(err)
			return
		}
		payload.EnginesVerdict = engines
	}()

	go func() {
		defer wg.Done()
		ipDomain, err := b.ipDomainModule.Build(ctx, vtData)
		if err != nil {
			setErr(err)
			return
		}
		payload.IPDomain = ipDomain
	}()

	go func() {
		defer wg.Done()
		summary, err := b.summaryModule.Build(ctx, vtData)
		if err != nil {
			setErr(err)
			return
		}
		payload.VTSummary = summary
	}()

	wg.Wait()
	if buildErr != nil {
		return models.ScanPayload{}, buildErr
	}

	return payload, nil
}
