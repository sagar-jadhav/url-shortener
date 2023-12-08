package metrics

import (
	"sort"
)

type Metric struct {
	Domain         string `json:"domain"`
	ShortenedCount int    `json:"count"`
}

type DomainMetrics struct {
	metrics map[string]int
}

func (m *DomainMetrics) Add(domain string) {
	if m.metrics == nil {
		m.metrics = make(map[string]int)
	}
	count, ok := m.metrics[domain]
	if !ok {
		m.metrics[domain] = 1
	} else {
		m.metrics[domain] = count + 1
	}
}

func (m *DomainMetrics) Get(domain string) int {
	return m.metrics[domain]
}

func (m *DomainMetrics) Sort() []Metric {
	metricsSlice := make([]Metric, 0)

	// prepare metrics slice from the map
	for domain, count := range m.metrics {
		metricsSlice = append(metricsSlice, Metric{
			Domain:         domain,
			ShortenedCount: count,
		})
	}

	// sort the slice
	sort.Slice(metricsSlice, func(i, j int) bool {
		return metricsSlice[i].ShortenedCount > metricsSlice[j].ShortenedCount
	})
	return metricsSlice
}
