package metrics

type Count struct {
	value   float64
	sampled bool
}

func (c *Count) addSample(sample *MetricSample, timestamp int64) {
	c.value = c.value + sample.Value
	c.sampled = true
}

func (c *Count) flush(timestamp int64) ([]*Point, error) {
	value, sampled := c.value, c.sampled
	c.value, c.sampled = 0, false

	if !sampled {
		return []*Point{}, NoDataError{}
	}

	return []*Point{{value, timestamp}}, nil
}
