package scheduler

import (
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
)

type jobQueue struct {
	interval time.Duration
	checks   []check.Check
}

func (jq *jobQueue) add(c check.Check) {
	jq.checks = append(jq.checks, c)
}

func (jq *jobQueue) start(checksPipe chan<- check.Check) {
	t := time.NewTicker(jq.interval)

	go func() {
		for _ = range t.C {

			for _, job := range jq.checks {
				checksPipe <- job
			}
		}
	}()
}
