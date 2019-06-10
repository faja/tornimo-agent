package statsd

import (
	"fmt"
	"net"
	"strings"
)

type listener interface {
	listen()
	stop()
}

type udpListener struct {
	conn   net.PacketConn
	pool   *packetPool
	buffer *packetBuffer
}

func newUDPListener(port uint, outputChannel chan []*packet, pool *packetPool) (listener, error) {
	var conn net.PacketConn
	var err error
	// TODO config host and port
	var url string = fmt.Sprintf("127.0.0.1:%d", port)

	conn, err = net.ListenPacket("udp", url)

	// TODO conn.(*net.UDPConn).SetReadBuffer(rcvbuf)

	if err != nil {
		// TODO return nicer error
		return nil, err
	}

	// TODO config flush interval
	buffer := newPacketBuffer(statsd_flush_interval, statsd_packet_buffer_size, outputChannel)

	l := &udpListener{
		conn:   conn,
		pool:   pool,
		buffer: buffer,
	}
	// TODO: log
	return l, nil
}

func (l *udpListener) listen() {
	//TODO log
	for {
		packet := l.pool.Get()
		//TODO stats
		n, _, err := l.conn.ReadFrom(packet.buffer)

		if err != nil {
			// connection has been closed
			if strings.HasSuffix(err.Error(), " use of closed network connection") {
				return
			}

			// TODO log and stats
			continue
		}

		if n == len(packet.buffer) {
			// TODO log and stats
		}

		packet.messages = packet.buffer[:n]
		l.buffer.append(packet)
	}
}

func (l *udpListener) stop() {
	l.buffer.close()
	l.conn.Close()
}
