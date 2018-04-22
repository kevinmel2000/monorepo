package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	metrics map[string]prometheus.Collector
)

var (
	namespace string

	summaryMetrics = metrics{
		"database_instrument": newSummaryMetrics("database", "instrument", "database instrumentation", []string{"type", "name"}),
	}
)

// InitObservation init
func InitObservation(servicename string) {
	namespace = servicename
}

func registerMetrics(m metrics) error {
	for _, val := range m {
		err := prometheus.Register(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func newSummaryMetrics(namespace, metricsName, doc string, labels []string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Name:      metricsName,
		Help:      doc,
	}, labels)
}
