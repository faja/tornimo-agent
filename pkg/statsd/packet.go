package statsd

type packet struct {
	buffer   []byte // underlying buffer, allocated only once
	messages []byte // actual message
}
