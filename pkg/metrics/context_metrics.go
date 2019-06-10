package metrics

import (
	"fmt"
	"log"
)

// TODO add indexing to ContextKey
type ContextKey string
type ContextMetrics map[ContextKey]Metric

// TODO pointer?
func (m ContextMetrics) AddSample(contextKey ContextKey, sample *MetricSample, timestamp int64) error {
	// TODO chek value if NaN

	if _, ok := m[contextKey]; !ok {
		switch sample.Mtype {
		case GaugeType:
			m[contextKey] = &Gauge{}
		case CountType:
			m[contextKey] = &Count{}
		default:
			return fmt.Errorf("unknown sample metric type: %v", sample.Mtype)
		}
	}

	m[contextKey].addSample(sample, timestamp)
	return nil
}

// TODO pointer?
func (m ContextMetrics) Flush(timestamp int64) []string {
	// TODO
	tornimoToken := "b04593ed-426c-425d-90d0-2c43c3f1576b"
	// TODO fix below
	metrics := make([]string, 0, 0)

	for key, metric := range m {
		points, err := metric.flush(timestamp)

		//TODO
		//if err != nil
		if err != nil {
			log.Printf("no data error while flushing metric '%s'\n", key)
		}

		for _, point := range points {
			// todo we dont want this for statsdSampler
			hostname := "mc.laptop"
			metricLine := fmt.Sprintf("%s.%s.%s %f %d\n", tornimoToken, hostname, key, point.Value, timestamp)
			metrics = append(metrics, metricLine)
		}
	}
	return metrics
}
