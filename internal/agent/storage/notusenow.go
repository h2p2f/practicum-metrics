package storage

import "fmt"

// Deprecated:URLMetrics is a method of the MetricStorage structure that generates a slice of
// URL addresses to send metrics to the server.
// Required for backward compatibility. Currently not used.
func (m *MetricStorage) URLMetrics(host string) []string {

	m.mut.Lock()
	defer m.mut.Unlock()

	var urls []string

	for metric, value := range m.gauge {
		generatedURL := fmt.Sprintf("%s/update/gauge/%s/%f", host, metric, value)
		urls = append(urls, generatedURL)
	}
	for metric, value := range m.counter {
		generatedURL := fmt.Sprintf("%s/update/counter/%s/%d", host, metric, value)
		urls = append(urls, generatedURL)
	}
	m.counter["PollCount"] = 0
	return urls
}
