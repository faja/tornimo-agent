package scheduler

import (
	"log"
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
)

type scheduler struct {
	checksPipe chan<- check.Check          // channel via which scheduler sends checks to runners
	enter      chan check.Check            // register checks to be scheduled
	jobQueues  map[time.Duration]*jobQueue // each interval has its own jobQueue
}

func NewScheduler(ch chan<- check.Check) *scheduler {
	log.Printf("[scheduler] creating and starting new scheduler\n")
	s := &scheduler{
		checksPipe: ch,
		enter:      make(chan check.Check),
		jobQueues:  make(map[time.Duration]*jobQueue),
	}

	go func() {
		for {
			select {
			case check := <-s.enter:
				s.scheduleCheck(check)
			}
		}
	}()

	return s
}

func (s *scheduler) scheduleCheck(c check.Check) {
	interval := c.Interval()
	jq, ok := s.jobQueues[interval]

	if !ok {
		log.Printf("[scheduler] creating new jobQueue[interval=%d]\n", int(interval.Truncate(time.Second).Seconds()))

		jq = &jobQueue{
			interval: interval,
			checks:   make([]check.Check, 0),
		}
		jq.start(s.checksPipe)
	}

	jq.add(c)
}

func (s *scheduler) GetEnterChan() chan<- check.Check {
	return s.enter
}
