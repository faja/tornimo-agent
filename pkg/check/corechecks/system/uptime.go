package system

import (
	"log"
	"time"

	"github.com/faja/tornimo-agent/pkg/aggregator"
	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/shirou/gopsutil/host"
)

type uptimeCheck struct {
	id check.ID
}

func NewUptimeCheck() *uptimeCheck {
	return &uptimeCheck{
		id: "uptimeCheck",
	}
}

func (c *uptimeCheck) Run() error {
	uptime, err := host.Uptime()
	if err != nil {
		log.Println(err)
		return err
	}
	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		log.Println(err)
		return err
	}

	sender.Gauge("system.uptime", float64(uptime), "", nil)
	sender.Commit()
	return nil
}

func (*uptimeCheck) Stop() error {
	return nil
}

func (*uptimeCheck) Interval() time.Duration {
	return time.Duration(time.Second * 15)
}

func (c *uptimeCheck) ID() check.ID {
	return c.id
}
