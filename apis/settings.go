package apis

import (
	"time"
)

type Settings struct {
	Timeout               time.Duration
	OperationDelayMethods []string
	OperationDelay        [2]time.Duration
	FastRel2Abs           bool
}

func DefaultSettings() *Settings {
	return &Settings{
		Timeout:               20 * time.Second,
		OperationDelayMethods: []string{"click", "swipe"},
		OperationDelay:        [2]time.Duration{200 * time.Microsecond, 200 * time.Microsecond},
		FastRel2Abs:           true,
	}
}

func (s *Settings) ImplicitlyWait(to time.Duration) {
	s.Timeout = to
}

func (s *Settings) operation_delay(operation_name string) func() {
	methodsContains := false
	for _, m := range s.OperationDelayMethods {
		if m == operation_name {
			methodsContains = true
			break
		}
	}
	if !methodsContains {
		return func() {}
	}
	before, after := s.OperationDelay[0], s.OperationDelay[1]
	time.Sleep(before)
	return func() {
		time.Sleep(after)
	}
}
