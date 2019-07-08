package system

import (
	"log"
	"time"

	"github.com/faja/tornimo-agent/pkg/aggregator"
	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/shirou/gopsutil/cpu"
)

type cpuCheck struct {
	id       check.ID
	cpuCount float64
	last     cpu.TimesStat
}

func NewCpuCheck() *cpuCheck {
	count, err := cpu.Counts(false)
	if err != nil {
		count = 1
	}

	return &cpuCheck{
		id:       "cpuCheck",
		cpuCount: float64(count),
	}
}

func (c *cpuCheck) Run() error {
	times, _ := cpu.Times(false)
	t := times[0]

	if c.last.Total() == 0 {
		c.last = t
		return nil
	}

	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		log.Println(err)
		return err
	}

	toPercent := 100 / (t.Total() - c.last.Total())

	user := (t.User - c.last.User) * toPercent
	system := (t.System - c.last.System) * toPercent
	idle := (t.Idle - c.last.Idle) * toPercent
	nice := (t.Nice - c.last.Nice) * toPercent
	iowait := (t.Iowait - c.last.Iowait) * toPercent
	irq := (t.Irq - c.last.Irq) * toPercent
	softirq := (t.Softirq - c.last.Softirq) * toPercent
	steal := (t.Steal - c.last.Steal) * toPercent
	guest := (t.Guest - c.last.Guest) * toPercent
	guestNice := (t.GuestNice - c.last.GuestNice) * toPercent

	c.last = t

	sender.Gauge("system.cpu.count", c.cpuCount, "", nil)
	sender.Gauge("system.cpu.user", user, "", nil)
	sender.Gauge("system.cpu.system", system, "", nil)
	sender.Gauge("system.cpu.idle", idle, "", nil)
	sender.Gauge("system.cpu.nice", nice, "", nil)
	sender.Gauge("system.cpu.iowait", iowait, "", nil)
	sender.Gauge("system.cpu.irq", irq, "", nil)
	sender.Gauge("system.cpu.softirq", softirq, "", nil)
	sender.Gauge("system.cpu.steal", steal, "", nil)
	sender.Gauge("system.cpu.guest", guest, "", nil)
	sender.Gauge("system.cpu.guestNice", guestNice, "", nil)
	sender.Commit()
	return nil
}

func (*cpuCheck) Stop() error {
	return nil
}

func (*cpuCheck) Interval() time.Duration {
	return time.Duration(time.Second * 15)
}

func (c *cpuCheck) ID() check.ID {
	return c.id
}
