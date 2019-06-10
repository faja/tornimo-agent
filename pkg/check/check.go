package check

import (
	"time"
)

type Check interface {
	ID() ID
	Interval() time.Duration
	Run() error
	Stop() error
}

type ID string
