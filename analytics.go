package logar

import (
	"sort"
	"time"

	"github.com/Lexographics/logar/models"
)

type AnalyticsSummary struct {
	TotalVisits      int64              `json:"total_visits"`       // Total number of requests received.
	UniqueVisitors   int64              `json:"unique_visitors"`    // Number of unique visitor IDs.
	ActiveVisitors   int64              `json:"active_visitors"`    // Number of active visitors (last 5 minutes)
	ErrorRate        float64            `json:"error_rate"`         // Rate of requests that are errors (e.g., 0.05 for 5%).
	AverageLatencyMs float64            `json:"average_latency_ms"` // Average request latency in milliseconds.
	P95LatencyMs     int64              `json:"p95_latency_ms"`     // 95th percentile request latency in milliseconds.
	P99LatencyMs     int64              `json:"p99_latency_ms"`     // 99th percentile request latency in milliseconds.
	TotalBytesSent   int64              `json:"total_bytes_sent"`   // Total bytes sent
	TotalBytesRecv   int64              `json:"total_bytes_recv"`   // Total bytes received
	TopPages         []PageStats        `json:"top_pages"`          // Top 5 most visited pages
	OSUsage          map[string]float64 `json:"os_usage"`           // OS usage distribution
	BrowserUsage     map[string]float64 `json:"browser_usage"`      // Browser usage distribution
	RefererUsage     map[string]float64 `json:"referer_usage"`      // Referer distribution
	InstanceStats    map[string]float64 `json:"instance_stats"`     // Request distribution by instance
}

// Statistics for a single page
type PageStats struct {
	Path       string  `json:"path"`
	Visits     int64   `json:"visits"`
	Percentage float64 `json:"percentage"`
}

type Analytics interface {
	// ID is automatically generated.
	RegisterRequest(log models.RequestLog) error

	GetStatistics(startTime time.Time, endTime time.Time) (AnalyticsSummary, error)
}

type AnalyticsImpl struct {
	core *AppImpl
}

func (a *AnalyticsImpl) GetApp() App {
	return a.core
}

func (a *AnalyticsImpl) RegisterRequest(log models.RequestLog) error {
	result := a.core.db.Create(&log)
	return result.Error
}

func (a *AnalyticsImpl) GetStatistics(startTime time.Time, endTime time.Time) (AnalyticsSummary, error) {
	summary := AnalyticsSummary{
		OSUsage:       map[string]float64{},
		BrowserUsage:  map[string]float64{},
		RefererUsage:  map[string]float64{},
		InstanceStats: map[string]float64{},
	}

	var totalVisits int64
	if err := a.core.db.Model(&models.RequestLog{}).Where("timestamp BETWEEN ? AND ?", startTime, endTime).Count(&totalVisits).Error; err != nil {
		return summary, err
	}
	summary.TotalVisits = totalVisits

	if totalVisits == 0 {
		return summary, nil
	}

	var uniqueVisitors int64
	if err := a.core.db.Model(&models.RequestLog{}).Where("timestamp BETWEEN ? AND ?", startTime, endTime).Distinct("visitor_id").Count(&uniqueVisitors).Error; err != nil {
		return summary, err
	}
	summary.UniqueVisitors = uniqueVisitors

	var errorCount int64
	if err := a.core.db.Model(&models.RequestLog{}).Where("timestamp BETWEEN ? AND ?", startTime, endTime).Where("status_code >= ?", 400).Count(&errorCount).Error; err != nil {
		return summary, err
	}
	if totalVisits > 0 {
		summary.ErrorRate = float64(errorCount) / float64(totalVisits)
	}

	var allLatenciesNano []int64
	if err := a.core.db.Model(&models.RequestLog{}).Where("timestamp BETWEEN ? AND ?", startTime, endTime).Pluck("latency", &allLatenciesNano).Error; err != nil {
		return summary, err
	}

	numLatencyRecords := len(allLatenciesNano)
	if numLatencyRecords == 0 {
		return summary, nil
	}

	var totalLatencySumNs int64
	latenciesMs := make([]int64, numLatencyRecords)
	for i, nanoSec := range allLatenciesNano {
		totalLatencySumNs += nanoSec
		latenciesMs[i] = nanoSec / int64(time.Millisecond)
	}

	summary.AverageLatencyMs = float64(totalLatencySumNs) / float64(numLatencyRecords) / float64(time.Millisecond)

	sort.Slice(latenciesMs, func(i, j int) bool {
		return latenciesMs[i] < latenciesMs[j]
	})

	p95Index := (numLatencyRecords * 95) / 100
	if p95Index >= numLatencyRecords {
		p95Index = numLatencyRecords - 1
	}
	if numLatencyRecords > 0 {
		summary.P95LatencyMs = latenciesMs[p95Index]
	}

	p99Index := (numLatencyRecords * 99) / 100
	if p99Index >= numLatencyRecords {
		p99Index = numLatencyRecords - 1
	}
	if numLatencyRecords > 0 {
		summary.P99LatencyMs = latenciesMs[p99Index]
	}

	var activeVisitorCount int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	err := a.core.db.Model(&models.RequestLog{}).Where("timestamp BETWEEN ? AND ?", startTime, endTime).Where("timestamp > ?", fiveMinutesAgo).Distinct("visitor_id").Count(&activeVisitorCount).Error
	if err != nil {
		return summary, err
	}
	summary.ActiveVisitors = activeVisitorCount

	type PageCount struct {
		Path  string
		Count int64
	}
	var pageCounts []PageCount
	if err := a.core.db.Model(&models.RequestLog{}).
		Select("path, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("path").
		Order("count desc").
		Limit(5).
		Scan(&pageCounts).Error; err != nil {
		return summary, err
	}

	summary.TopPages = make([]PageStats, len(pageCounts))
	for i, pc := range pageCounts {
		summary.TopPages[i] = PageStats{
			Path:       pc.Path,
			Visits:     pc.Count,
			Percentage: float64(pc.Count) / float64(totalVisits) * 100,
		}
	}

	type OSCount struct {
		OS    string
		Count int64
	}
	var osCounts []OSCount
	if err := a.core.db.Model(&models.RequestLog{}).
		Select("os, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("os").
		Scan(&osCounts).Error; err != nil {
		return summary, err
	}

	for _, oc := range osCounts {
		summary.OSUsage[oc.OS] = float64(oc.Count) / float64(totalVisits) * 100
	}

	type BrowserCount struct {
		Browser string
		Count   int64
	}
	var browserCounts []BrowserCount
	if err := a.core.db.Model(&models.RequestLog{}).
		Select("browser, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("browser").
		Scan(&browserCounts).Error; err != nil {
		return summary, err
	}

	for _, bc := range browserCounts {
		summary.BrowserUsage[bc.Browser] = float64(bc.Count) / float64(totalVisits) * 100
	}

	var totalBytesSent, totalBytesRecv int64
	if err := a.core.db.Model(&models.RequestLog{}).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Select("COALESCE(SUM(bytes_sent), 0), COALESCE(SUM(bytes_recv), 0)").
		Row().
		Scan(&totalBytesSent, &totalBytesRecv); err != nil {
		return summary, err
	}
	summary.TotalBytesSent = totalBytesSent
	summary.TotalBytesRecv = totalBytesRecv

	type InstanceCount struct {
		Instance string
		Count    int64
	}
	var instanceCounts []InstanceCount
	if err := a.core.db.Model(&models.RequestLog{}).
		Select("instance, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("instance").
		Scan(&instanceCounts).Error; err != nil {
		return summary, err
	}

	for _, ic := range instanceCounts {
		summary.InstanceStats[ic.Instance] = float64(ic.Count) / float64(totalVisits) * 100
	}

	type RefererCount struct {
		Referer string
		Count   int64
	}
	var refererCounts []RefererCount
	if err := a.core.db.Model(&models.RequestLog{}).
		Select("referer, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("referer").
		Scan(&refererCounts).Error; err != nil {
		return summary, err
	}

	for _, rc := range refererCounts {
		summary.RefererUsage[rc.Referer] = float64(rc.Count) / float64(totalVisits) * 100
	}

	return summary, nil
}
