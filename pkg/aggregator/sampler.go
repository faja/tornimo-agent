package aggregator

import (
	"bytes"
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
func (s *sampler) flush(timestamp int64) {

	// TODO add error handling
	// TODO
	var b bytes.Buffer

	// TODO send this to serializer
	// for now there is only 0 bucket
	for _, cm := range s.bucket {
		// TODO add []Serie
		metrics := cm.Flush(timestamp)

		for _, metricLine := range metrics {

			_, err := b.WriteString(metricLine)
			if err != nil {
				log.Print(err)
			}
		}

	}
	aggregatorInstance.forwarder.SubmitSeries(b.Bytes())
}
