package collector

import (
	"log"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/check/corechecks"
	"github.com/faja/tornimo-agent/pkg/collector/runner"
	"github.com/faja/tornimo-agent/pkg/collector/scheduler"
)

func NewCollector() {
	log.Printf("[collector] creating and starting new collector\n")
	checkChan := make(chan check.Check)

	_ = runner.NewRunner(checkChan)
	s := scheduler.NewScheduler(checkChan)

	corechecks.LoadChecks(s.GetEnterChan())
}
