package statsd

import "sync"

type packetPool struct {
	pool sync.Pool
}

// TODO: write documentation about reusing byte slices

func newPacketPool(buffer uint) *packetPool {
	return &packetPool{
		pool: sync.Pool{
			New: func() interface{} {
				// allocate big buffer once, by default 16k
				// after that we manage messages that point to underlying buffer
				packet := &packet{
					buffer: make([]byte, buffer),
				}
				packet.messages = packet.buffer[0:0]
				return packet
			},
		},
	}
}

func (p *packetPool) Get() *packet {
	return p.pool.Get().(*packet)
}

func (p *packetPool) Put(packet *packet) {
	p.pool.Put(packet)
}
