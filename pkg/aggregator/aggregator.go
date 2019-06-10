package aggregator

import (
	"sync"
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/forwarder"
	"github.com/faja/tornimo-agent/pkg/metrics"
)

const (
	tornimoToken string = "b04593ed-426c-425d-90d0-2c43c3f1576b"
)

var (
	aggregatorInstance *Aggregator
	aggregatorInit     sync.Once
)

type Aggregator struct {
	metricIn         chan senderMetricSample
	bufferedMetricIn chan []*metrics.MetricSample
	tickerChan       <-chan time.Time
	m                sync.Mutex
	sampler          map[check.ID]*sampler
	statsdSampler    *sampler
	forwarder        forwarder.Forwarder
}

func InitAggregator(f forwarder.Forwarder) *Aggregator {
	aggregatorInit.Do(func() {
		aggregatorInstance = newAggregator(f)
	})

	return aggregatorInstance
}

func newAggregator(forwarder forwarder.Forwarder) *Aggregator {
	agg := &Aggregator{
		metricIn:         make(chan senderMetricSample, 100),
		bufferedMetricIn: make(chan []*metrics.MetricSample, 100),
		tickerChan:       time.NewTicker(time.Second * 15).C,
		sampler:          make(map[check.ID]*sampler),
		// TODO
		statsdSampler: newSampler(),
		forwarder:     forwarder,
	}
	go agg.run()

	return agg
}

func (agg *Aggregator) run() {
	for {
		select {
		case <-agg.tickerChan:
			now := time.Now().Unix()

			// check samplers
			// TODO send this to serializer
			for _, v := range agg.sampler {
				v.flush(now)
			}

			// statsd sampler
			// TODO send this to serializer
			agg.statsdSampler.flush(now)

		case m := <-agg.metricIn:
			// TODO for now lets just skip COMMIT()
			if m.commit {
				continue
			}
			sampler, ok := agg.sampler[m.id]
			if !ok {
				sampler = newSampler()
				agg.sampler[m.id] = sampler
			}

			sampler.addSample(m.metric, time.Now().Unix())

		case samples := <-agg.bufferedMetricIn:
			now := time.Now().Unix()
			for _, sample := range samples {
				agg.statsdSampler.addSample(sample, now)
			}
		}
	}
}

// TODO write it nicer
func (agg *Aggregator) GetMetricsChan() chan []*metrics.MetricSample {
	return agg.bufferedMetricIn
}
