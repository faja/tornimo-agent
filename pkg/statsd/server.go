package statsd

import (
	"time"

	"github.com/faja/tornimo-agent/pkg/metrics"
)

// TODO config
var statsd_queue_size uint = 100
var statsd_packet_buffer_size uint = 512
var statsd_buffer_size uint = 1024 * 16 // 16k udp read buffer
var statsd_flush_interval time.Duration = time.Duration(time.Second * 15)

type Server struct {
	listeners   []listener
	packetsChan chan []*packet
	stopChan    chan bool
	pool        *packetPool

	// TODO
	/*
				Statistics            *util.Stats
		    Started               bool
		    health                *health.Handle
		    metricPrefix          string
		    metricPrefixBlacklist []string
		    defaultHostname       string
		    histToDist            bool
		    histToDistPrefix      string
		    extraTags             []string
		    debugMetricsStats     bool
		    metricsStats          map[string]metricStat
		    statsLock             sync.Mutex
	*/
}

func NewServer(port uint, metricOut chan<- []*metrics.MetricSample) (*Server, error) {
	packetsChan := make(chan []*packet, statsd_queue_size)
	pool := newPacketPool(statsd_buffer_size)

	l, err := newUDPListener(port, packetsChan, pool)
	if err != nil {
		// TODO: nicer err
		return nil, err
	}

	listeners := make([]listener, 1)
	listeners[0] = l

	s := &Server{
		listeners:   listeners,
		packetsChan: packetsChan,
		stopChan:    make(chan bool),
		pool:        pool,
	}

	s.start(metricOut)

	return s, nil
}

func (s *Server) Stop() {
	// TODO lock?
	close(s.stopChan)
	for _, l := range s.listeners {
		l.stop()
	}
	// TODO internalState?
}

func (s *Server) start(metricOut chan<- []*metrics.MetricSample) {
	for _, l := range s.listeners {
		go l.listen()
	}

	// TODO config number of workers?
	workers := 2
	for i := 0; i < workers; i++ {
		go s.work(metricOut)
	}
}

func (s *Server) work(metricOut chan<- []*metrics.MetricSample) {
	for {
		select {
		case <-s.stopChan:
			return
		case packets := <-s.packetsChan:
			// TODO stats and log
			metricSamples := make([]*metrics.MetricSample, 0, len(packets))

			for _, packet := range packets {
				//TODO stats and log
				metricSamples = parsePacket(packet, metricSamples)
				// TODO stats and log
				s.pool.Put(packet)
			}

			// TODO stats and log
			if len(metricSamples) > 0 {
				metricOut <- metricSamples
			}
		}
	}
}
