package forwarder

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/faja/tornimo-agent/pkg/state"
)

/*
 * grafana: https://tornimo.tornimo.io
 * put:     put.tornimo.tornimo.io 2003
 * token:   b04593ed-426c-425d-90d0-2c43c3f1576b
 */

var tornimoAddress = "put.tornimo.tornimo.io:2003"

type worker struct {
	id            uint32
	internalState state.State
	conn          net.Conn

	highPrioChan <-chan transaction
	lowPrioChan  <-chan transaction
	requeueChan  chan<- transaction

	stopChan    chan bool
	stoppedChan chan bool

	m sync.Mutex // To control Start/Stop races
}

func newWorker(id uint32, highPrioChan <-chan transaction, lowPrioChan <-chan transaction, requeueChan chan<- transaction) (*worker, error) {
	log.Printf("[OK] creating worker [id:%d]\n", id)

	c, err := net.Dial("tcp", tornimoAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not create forwarder.worker: %v\n", err)
	}

	log.Printf("[OK] worker [id:%d] created\n", id)
	return &worker{
		id:            id,
		internalState: state.Stopped,
		conn:          c,
		highPrioChan:  highPrioChan,
		lowPrioChan:   lowPrioChan,
		requeueChan:   requeueChan,
		stopChan:      make(chan bool),
		stoppedChan:   make(chan bool),
	}, nil
}

func (w *worker) stop() {
	log.Printf("[OK] stopping worker [id:%d]\n", w.id)

	w.m.Lock()
	defer w.m.Unlock()

	if w.internalState == state.Stopped {
		log.Printf("[WARN] worker [id %d] is already stopped\n", w.id)
		return
	}

	w.stopChan <- true
	<-w.stoppedChan
	w.internalState = state.Stopped

	log.Printf("[OK] worker [id:%d] stopped\n", w.id)
}

func (w *worker) start() {
	log.Printf("[OK] starting worker [id:%d]\n", w.id)

	w.m.Lock()
	defer w.m.Unlock()

	if w.internalState == state.Started {
		log.Printf("[WARN] worker [id %d] is already started\n", w.id)
		return
	}

	/*
	 * IMPORTANT:
	 * this function should return only by sending a signal via stopChan
	 * otherwise call to stop() would block
	 */
	go func() {

		defer close(w.stoppedChan)

		for {
			select {
			case <-w.stopChan:
				log.Printf("worker main loop got STOP signal")
				return
			case t := <-w.highPrioChan:
				err := t.process(context.Background(), w.conn)
				if err != nil {
					log.Println("[worker] - could nod send data to tornimo:", err)
				}
			}
		}
	}()

	w.internalState = state.Started
	log.Printf("[OK] worker [id:%d] started\n", w.id)
}
