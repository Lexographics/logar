package logar

import (
	"fmt"
	"time"
)

type Timer struct {
	start  time.Time
	logger *AppImpl
}

func (l *AppImpl) NewTimer() *Timer {
	return &Timer{
		start:  time.Now(),
		logger: l,
	}
}

func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

func (t *Timer) Log(model Model, message string, category string) error {
	return t.logger.GetLogger().Log(
		model,
		fmt.Sprintf("\"%s\" took %s", message, t.Elapsed().String()),
		category,
	)
}

func (t *Timer) Reset() {
	t.start = time.Now()
}

func (t *Timer) StartTime() time.Time {
	return t.start
}
