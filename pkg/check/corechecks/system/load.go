package system

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/load"

	"github.com/faja/tornimo-agent/pkg/aggregator"
	"github.com/faja/tornimo-agent/pkg/check"
)

type loadCheck struct {
	id check.ID
}

func NewLoadCheck() *loadCheck {
	return &loadCheck{
		id: "loadCheck",
	}
}

func (c *loadCheck) Run() error {
	avg, err := load.Avg()
	if err != nil {
		log.Println(err)
	}
	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		log.Println(err)
		return err
	}

	sender.Gauge("system.loadavg.load1", avg.Load1, "", nil)
	sender.Gauge("system.loadavg.load5", avg.Load5, "", nil)
	sender.Gauge("system.loadavg.load15", avg.Load15, "", nil)
	sender.Commit()
	return nil
}

func (*loadCheck) Stop() error {
	return nil
}

func (*loadCheck) Interval() time.Duration {
	return time.Duration(time.Second * 15)
}

func (c *loadCheck) ID() check.ID {
	return c.id
}
