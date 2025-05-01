package logar

import (
	"fmt"
	"time"
)

type Timer struct {
	start  time.Time
	logger *Logger
}

func (l *Logger) NewTimer() *Timer {
	return &Timer{
		start:  time.Now(),
		logger: l,
	}
}

func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

func (t *Timer) Log(model string, message string, category string) error {
	return t.logger.Log(
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
