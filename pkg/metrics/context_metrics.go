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
func (m ContextMetrics) Flush(timestamp int64) []*Serie {
	// TODO fix below
	series := make([]*Serie, 0, 0)

	for metricName, metricPoints := range m {
		points, err := metricPoints.flush(timestamp)
		//TODO better error handling
		if err != nil {
			log.Printf("no data error while flushing metric '%s'\n", metricName)
		}

		serie := &Serie{
			Name:   string(metricName),
			Points: points,
		}
		series = append(series, serie)
	}

	return series
}
