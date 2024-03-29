package forwarder

import (
	"fmt"
	"log"
	"sync"

	"github.com/faja/tornimo-agent/pkg/state"
)

var chanBufferSize uint8 = 100

type Forwarder interface {
	SubmitSeries(b []byte)
}

type defaultForwarder struct {
	tornimoAddress  string
	numberOfWorkers uint32
	internalState   state.State

	highPrio            chan transaction
	lowPrio             chan transaction
	requeuedTransaction chan transaction

	workers []*worker
	m       sync.Mutex
}

func NewDefaultForwarder(tornimoAddress string) *defaultForwarder {
	return &defaultForwarder{
		tornimoAddress:      tornimoAddress,
		numberOfWorkers:     2,
		internalState:       state.Stopped,
		highPrio:            make(chan transaction),
		lowPrio:             make(chan transaction),
		requeuedTransaction: make(chan transaction),
		workers:             []*worker{},
	}
}

func (f *defaultForwarder) Start() error {
	log.Printf("[forwarder] starting forwarder\n")
	f.m.Lock()
	defer f.m.Unlock()

	if f.internalState == state.Started {
		// TODO logging
		log.Printf("[forwarder][ERROR] defaultForwarder is already started\n")
		return fmt.Errorf("defaultForwarder is already started\n")
	}

	for id := uint32(1); id <= f.numberOfWorkers; id++ {
		w, err := newWorker(id, f.tornimoAddress, f.highPrio, f.lowPrio, f.requeuedTransaction)
		if err != nil {
			return err
		}
		f.workers = append(f.workers, w)
	}

	for _, w := range f.workers {
		w.start()
	}

	f.internalState = state.Started
	// TODO logging
	// log.Printf("[forwarder] defaultForwarder started\n")
	return nil
}

func (f *defaultForwarder) Stop() {
	// TODO logging
	//log.Printf("[forwarder] stopping defaultForwarder\n")
	f.m.Lock()
	defer f.m.Unlock()

	if f.internalState == state.Started {
		log.Printf("[forwarder] defaultForwarder is already stopped\n")
	}

	for _, w := range f.workers {
		w.stop()
	}

	f.internalState = state.Stopped
	log.Printf("[forwarder] defaultForwarder stopped\n")
}

func (f *defaultForwarder) SubmitSeries(b []byte) {
	f.highPrio <- &defaultTransaction{b}
}
