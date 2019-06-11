package serializer

import (
	"bytes"
	"fmt"

	"github.com/faja/tornimo-agent/pkg/forwarder"
	"github.com/faja/tornimo-agent/pkg/metrics"
)

type Serializer interface {
	SendSeries([]*metrics.Serie) error
}

type defaultSerializer struct {
	forwarder forwarder.Forwarder
	token     string
}

func NewSerializer(forwarder forwarder.Forwarder, token string) Serializer {
	return &defaultSerializer{
		forwarder: forwarder,
		token:     token,
	}
}

func (s *defaultSerializer) SendSeries(series []*metrics.Serie) error {
	var b bytes.Buffer
	for _, serie := range series {
		for _, point := range serie.Points {
			b.WriteString(serializeSerie(s.token, serie.Name, point.Value, point.Timestamp))
		}
	}

	// TODO add error to forwarder
	s.forwarder.SubmitSeries(b.Bytes())
	// TODO error
	return nil
}

func serializeSerie(token, serie string, value float64, timestamp int64) string {
	return fmt.Sprintf("%s.%s %f %d\n", token, serie, value, timestamp)
}
