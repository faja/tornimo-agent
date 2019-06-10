package metrics

type Point struct {
	Value     float64
	Timestamp int64
}

type NoDataError struct{}

func (e NoDataError) Error() string {
	// TODO
	return "No data"
}
