package runner

import (
	"log"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/state"
)

type runner struct {
	numberOfWorkers uint8
	internalState   state.State
	pending         <-chan check.Check
	// TODO make use of runningChecks
	//runningChecks   map[check.ID]check.Check
	// TODO make use of m
	//m sync.Mutex
}

func NewRunner(ch <-chan check.Check) *runner {
	log.Printf("[runner] starting runner\n")

	r := &runner{
		numberOfWorkers: 2,
		internalState:   state.Started,
		pending:         ch,
		//runningChecks:   make(map[check.ID]check.Check),
	}

	for i := uint8(0); i < r.numberOfWorkers; i++ {
		go r.work(i)
	}

	return r
}

func (r *runner) work(i uint8) {
	log.Printf("[runner] starting runner worker[id:%d]\n", i)
	for check := range r.pending {
		check.Run()
	}
}
