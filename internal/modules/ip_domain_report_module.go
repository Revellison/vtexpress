package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"vtgui/internal/models"
)

var (
	ipv4Pattern   = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	domainPattern = regexp.MustCompile(`\b(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}\b`)
)

type IPDomainReportModule struct {
	httpClient *http.Client
}

func NewIPDomainReportModule() *IPDomainReportModule {
	return &IPDomainReportModule{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *IPDomainReportModule) Build(ctx context.Context, vtData map[string]any) (models.IPDomainReport, error) {
	stringValues := collectStrings(vtData)
	ips := make(map[string]struct{})
	domains := make(map[string]struct{})

	for _, value := range stringValues {
		for _, token := range ipv4Pattern.FindAllString(value, -1) {
			parsed := net.ParseIP(token)
			if parsed == nil || parsed.To4() == nil {
				continue
			}
			ips[token] = struct{}{}
		}
		for _, token := range domainPattern.FindAllString(strings.ToLower(value), -1) {
			if net.ParseIP(token) != nil {
				continue
			}
			domains[token] = struct{}{}
		}
	}

	for domain := range domains {
		resolved, _ := resolveDomain(ctx, domain)
		if resolved != "" {
			ips[resolved] = struct{}{}
		}
	}

	ipList := make([]string, 0, len(ips))
	for ip := range ips {
		ipList = append(ipList, ip)
	}
	sort.Strings(ipList)

	ipDetails := make(map[string]models.IPInfo, len(ipList))
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, ip := range ipList {
		ipCopy := ip
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := models.IPInfo{Address: ipCopy}
			enriched, err := m.lookupIP(ctx, ipCopy)
			if err == nil {
				info = enriched
			}
			mu.Lock()
			ipDetails[ipCopy] = info
			mu.Unlock()
		}()
	}
	wg.Wait()

	result := models.IPDomainReport{
		IPs:     make([]models.IPInfo, 0, len(ipDetails)),
		Domains: make([]models.DomainInfo, 0, len(domains)),
	}

	for _, ip := range ipList {
		result.IPs = append(result.IPs, ipDetails[ip])
	}

	domainList := make([]string, 0, len(domains))
	for domain := range domains {
		domainList = append(domainList, domain)
	}
	sort.Strings(domainList)

	for _, domain := range domainList {
		resolved, _ := resolveDomain(ctx, domain)
		domainInfo := models.DomainInfo{Domain: domain, Resolved: resolved}
		if resolved != "" {
			if ipInfo, ok := ipDetails[resolved]; ok {
				domainInfo.ASN = ipInfo.ASN
				domainInfo.Provider = ipInfo.Provider
				domainInfo.Region = ipInfo.Region
			}
		}
		result.Domains = append(result.Domains, domainInfo)
	}

	return result, nil
}

func (m *IPDomainReportModule) lookupIP(ctx context.Context, ip string) (models.IPInfo, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipwho.is/"+ip, nil)
	if err != nil {
		return models.IPInfo{}, err
	}
	response, err := m.httpClient.Do(request)
	if err != nil {
		return models.IPInfo{}, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return models.IPInfo{}, fmt.Errorf("ipwhois status: %d", response.StatusCode)
	}

	var payload struct {
		Success bool   `json:"success"`
		IP      string `json:"ip"`
		Region  string `json:"region"`
		Country string `json:"country"`
		Conn    struct {
			ISP string `json:"isp"`
			ASN string `json:"asn"`
		} `json:"connection"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return models.IPInfo{}, err
	}
	if !payload.Success {
		return models.IPInfo{}, fmt.Errorf("ipwhois unsuccessful response")
	}

	region := strings.TrimSpace(strings.Join([]string{payload.Region, payload.Country}, ", "))
	return models.IPInfo{
		Address:  ip,
		ASN:      payload.Conn.ASN,
		Provider: payload.Conn.ISP,
		Region:   strings.Trim(region, ", "),
	}, nil
}

func resolveDomain(ctx context.Context, domain string) (string, error) {
	systemResolver := net.DefaultResolver
	ips, err := systemResolver.LookupIPAddr(ctx, domain)
	if err == nil {
		for _, ip := range ips {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	fallbackResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 3 * time.Second}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}

	fallbackIPs, fallbackErr := fallbackResolver.LookupIPAddr(ctx, domain)
	if fallbackErr != nil {
		return "", fallbackErr
	}
	for _, ip := range fallbackIPs {
		if ip.IP.To4() != nil {
			return ip.IP.String(), nil
		}
	}

	return "", fmt.Errorf("no ipv4 address resolved for domain: %s", domain)
}

func collectStrings(value any) []string {
	results := make([]string, 0)
	var walk func(v any)
	walk = func(v any) {
		switch t := v.(type) {
		case map[string]any:
			for _, vv := range t {
				walk(vv)
			}
		case []any:
			for _, vv := range t {
				walk(vv)
			}
		case string:
			if t != "" {
				results = append(results, t)
			}
		}
	}
	walk(value)
	return results
}
