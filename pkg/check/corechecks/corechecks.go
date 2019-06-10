package corechecks

import (
	"time"

	"github.com/faja/tornimo-agent/pkg/check"
	"github.com/faja/tornimo-agent/pkg/check/corechecks/system"
)

type BaseCheck struct {
	name     string
	id       check.ID
	interval time.Duration
}

func LoadChecks(enter chan<- check.Check) {
	cc := genCoreChecks()

	for _, check := range cc {
		enter <- check
	}
}

func genCoreChecks() []check.Check {
	coreChecks := make([]check.Check, 0)
	coreChecks = append(coreChecks, system.NewLoadCheck())
	return coreChecks
}