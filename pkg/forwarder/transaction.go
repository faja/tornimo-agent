package forwarder

import (
	"context"
	"net"
)

type transaction interface {
	process(context.Context, net.Conn) error
}

type defaultTransaction struct {
	data []byte
}

func (t *defaultTransaction) process(ctx context.Context, conn net.Conn) error {
	_, err := conn.Write(t.data)
	return err
}
