package metrics

type MetricType uint8

const (
	GaugeType MetricType = iota
	CountType
)

type MetricSample struct {
	Name      string
	Value     float64
	Tags      []string
	Timestamp int64
	Mtype     MetricType
}

type Metric interface {
	addSample(*MetricSample, int64)
	flush(int64) ([]*Point, error)
}
