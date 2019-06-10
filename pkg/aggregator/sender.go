package aggregator

import (
	"errors"
	"sync"
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/metrics"
)

var senderPool *checkSenderPool = &checkSenderPool{senders: make(map[check.ID]Sender)}

type Sender interface {
	Commit()
	Gauge(metric string, value float64, hostname string, tags []string)
}

type checkSenderPool struct {
	senders map[check.ID]Sender
	m       sync.Mutex
}

type checkSender struct {
	id              check.ID
	defaultHostname string
	smsOut          chan<- senderMetricSample
}

type senderMetricSample struct {
	id     check.ID
	metric *metrics.MetricSample
	commit bool
}

func newCheckSender(id check.ID, defaultHostname string) *checkSender {
	return &checkSender{
		id:              id,
		defaultHostname: defaultHostname,
		smsOut:          aggregatorInstance.metricIn,
	}
}

func GetDefaultSender() (Sender, error) {
	var defaultChecID check.ID // zero value empty string ""
	return GetSender(defaultChecID)
}

func GetSender(id check.ID) (Sender, error) {
	if aggregatorInstance == nil {
		return nil, errors.New("Aggregator was not initialized")
	}

	return senderPool.getSender(id)
}

func (sp *checkSenderPool) getSender(id check.ID) (Sender, error) {
	var sender Sender

	sp.m.Lock()
	sender, ok := sp.senders[id]
	if !ok {
		sender = newCheckSender(id, "")
		sp.senders[id] = sender
	}
	sp.m.Unlock()

	return sender, nil
}

func (s *checkSender) Commit() {
	s.smsOut <- senderMetricSample{s.id, &metrics.MetricSample{}, true}
}

func (s *checkSender) Gauge(metricName string, value float64, hostname string, tags []string) {
	s.sendMetric(metricName, value, hostname, tags, metrics.GaugeType)
}

func (s *checkSender) sendMetric(metricName string, value float64, hostname string, metricTags []string, metircType metrics.MetricType) {
	s.smsOut <- senderMetricSample{
		s.id,
		&metrics.MetricSample{
			Name:      metricName,
			Value:     value,
			Tags:      metricTags,
			Timestamp: time.Now().Unix(),
		},
		false,
	}
}
