package statsd

import (
	"sync"
	"time"
)

type packetBuffer struct {
	packets       []*packet
	flushTimer    *time.Ticker
	bufferSize    uint
	outputChannel chan []*packet
	closeChannel  chan struct{}
	m             sync.Mutex
}

func newPacketBuffer(flushInterval time.Duration, bufferSize uint, outputChannel chan []*packet) *packetBuffer {
	pb := &packetBuffer{
		packets:       make([]*packet, 0, bufferSize),
		flushTimer:    time.NewTicker(flushInterval),
		bufferSize:    bufferSize,
		outputChannel: outputChannel,
		closeChannel:  make(chan struct{}),
	}
	go pb.flushLoop()
	return pb
}

func (pb *packetBuffer) flushLoop() {
	for {
		select {
		case <-pb.flushTimer.C:
			pb.m.Lock()
			pb.flush()
			pb.m.Unlock()
		case <-pb.closeChannel:
			return
		}
	}
}

func (pb *packetBuffer) flush() {
	if len(pb.packets) > 0 {
		pb.outputChannel <- pb.packets
		pb.packets = make([]*packet, 0, pb.bufferSize)
	}
}

func (pb *packetBuffer) append(packet *packet) {
	pb.m.Lock()
	pb.packets = append(pb.packets, packet)
	if uint(len(pb.packets)) >= pb.bufferSize {
		pb.flush()
	}
	pb.m.Unlock()
}

func (pb *packetBuffer) close() {
	close(pb.closeChannel)
}
