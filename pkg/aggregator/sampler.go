package aggregator

import (
	"log"

	"github.com/faja/tornimo-agent/pkg/metrics"
)

//TODO: add contexmetric abstraction
type sampler struct {
	// TODO: refactor!
	bucket map[int64]metrics.ContextMetrics // time.now.unix
}

func newSampler() *sampler {
	return &sampler{
		bucket: map[int64]metrics.ContextMetrics{},
	}
}

func (s *sampler) addSample(metricSample *metrics.MetricSample, timestamp int64) {
	contextKey := metrics.ContextKey(metricSample.Name)

	// TODO: for now only 1 bucket "0"
	bucketMetrics, ok := s.bucket[0]
	if !ok {
		bucketMetrics = metrics.ContextMetrics{}
		s.bucket[0] = bucketMetrics
	}

	if err := bucketMetrics.AddSample(contextKey, metricSample, timestamp); err != nil {
		// TODO
		log.Print("Could not add sample '%s': %v", metricSample.Name, err)
	}
}

// TODO
func (s *sampler) flush(timestamp int64) []*metrics.Serie {
	// TODO add error handling

	series := make([]*metrics.Serie, 0, 0)

	// for now there is only 0 bucket
	for _, cm := range s.bucket {
		series = append(series, cm.Flush(timestamp)...)
	}

	return series
}
