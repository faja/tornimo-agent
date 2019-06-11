package aggregator

import (
	"log"
	"sync"
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/metrics"
	"github.com/faja/tornimo-agent/pkg/serializer"
)

var (
	aggregatorInstance *Aggregator
	aggregatorInit     sync.Once
)

type Aggregator struct {
	defaultHostname  string
	metricIn         chan senderMetricSample
	bufferedMetricIn chan []*metrics.MetricSample
	tickerChan       <-chan time.Time
	samplers         map[check.ID]*sampler
	statsdSampler    *sampler
	serializer       serializer.Serializer
	m                sync.Mutex
}

func InitAggregator(defaultHostname string, s serializer.Serializer) *Aggregator {
	aggregatorInit.Do(func() {
		aggregatorInstance = newAggregator(defaultHostname, s)
	})

	return aggregatorInstance
}

func newAggregator(defaultHostname string, s serializer.Serializer) *Aggregator {
	agg := &Aggregator{
		defaultHostname:  defaultHostname,
		metricIn:         make(chan senderMetricSample, 100),
		bufferedMetricIn: make(chan []*metrics.MetricSample, 100),
		tickerChan:       time.NewTicker(time.Second * 15).C,
		samplers:         make(map[check.ID]*sampler),
		statsdSampler:    newSampler(),
		serializer:       s,
	}

	// TODO not sure if this is good idea
	// lets put statsdSampler as a default one
	var defaultChecID check.ID // zero value empty string ""
	agg.samplers[defaultChecID] = agg.statsdSampler

	go agg.run()

	return agg
}

func (agg *Aggregator) flush() {
	// TODO
	// for now all the series gets timestampped at this point
	now := time.Now().Unix()
	series := make([]*metrics.Serie, 0, 0)

	// fetch all series from all the samplers
	for _, sampler := range agg.samplers {
		series = append(series, sampler.flush(now)...)
	}

	go func() {
		log.Printf("[aggregator] flushing %d series\n", len(series))
		err := agg.serializer.SendSeries(series)
		// TODO handle error
		_ = err
	}()
}

func (agg *Aggregator) run() {
	for {
		select {
		case <-agg.tickerChan:
			// flush all the samples to serializer
			agg.flush()

		case m := <-agg.metricIn:
			// TODO for now lets just skip COMMIT()
			if m.commit {
				continue
			}
			sampler, ok := agg.samplers[m.id]
			if !ok {
				sampler = newSampler()
				agg.samplers[m.id] = sampler
			}

			sampler.addSample(m.metric, time.Now().Unix())

		case samples := <-agg.bufferedMetricIn:
			// TODO: what to do with now? should I pass it to addSample?
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
