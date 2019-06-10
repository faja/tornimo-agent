package metrics

type Gauge struct {
	value   float64
	sampled bool
}

func (g *Gauge) addSample(sample *MetricSample, timestamp int64) {
	g.value = sample.Value
	g.sampled = true
}

func (g *Gauge) flush(timestamp int64) ([]*Point, error) {
	value, sampled := g.value, g.sampled
	g.value, g.sampled = 0, false

	if !sampled {
		return []*Point{}, NoDataError{}
	}

	return []*Point{{value, timestamp}}, nil
}
