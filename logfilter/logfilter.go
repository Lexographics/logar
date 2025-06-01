package logfilter

import (
	"strings"
	"time"

	"sadk.dev/logar/models"
)

type Filter struct {
	conditions []Condition
}

type Condition func(log models.Log) bool

func NewFilter(conditions ...Condition) Filter {
	return Filter{
		conditions: conditions,
	}
}

func (f *Filter) Evaluate(log models.Log) bool {
	for _, c := range f.conditions {
		if !c(log) {
			return false
		}
	}
	return true
}

func MessageContains(substr string) Condition {
	return func(log models.Log) bool {
		return strings.Contains(log.Message, substr)
	}
}

func CategoryContains(substr string) Condition {
	return func(log models.Log) bool {
		return strings.Contains(log.Category, substr)
	}
}

func ModelContains(substr string) Condition {
	return func(log models.Log) bool {
		return strings.Contains(string(log.Model), substr)
	}
}

func IsModel(model models.Model) Condition {
	return func(log models.Log) bool {
		return log.Model == model
	}
}

func IsCategory(category string) Condition {
	return func(log models.Log) bool {
		return log.Category == category
	}
}

func IsSeverity(severity models.Severity) Condition {
	return func(log models.Log) bool {
		return log.Severity == severity
	}
}

func IsSeverityAtLeast(severity models.Severity) Condition {
	return func(log models.Log) bool {
		return log.Severity >= severity
	}
}

func IsSeverityAtMost(severity models.Severity) Condition {
	return func(log models.Log) bool {
		return log.Severity <= severity
	}
}

func IsSeverityBetween(min models.Severity, max models.Severity) Condition {
	return func(log models.Log) bool {
		return log.Severity >= min && log.Severity <= max
	}
}

func TimeBetween(min, max time.Time) Condition {
	return func(log models.Log) bool {
		return log.CreatedAt.After(min) && log.CreatedAt.Before(max)
	}
}

func TimeAfter(min time.Time) Condition {
	return func(log models.Log) bool {
		return log.CreatedAt.After(min)
	}
}

func TimeBefore(max time.Time) Condition {
	return func(log models.Log) bool {
		return log.CreatedAt.Before(max)
	}
}

func DayOfWeekIn(days ...time.Weekday) Condition {
	return func(log models.Log) bool {
		for _, d := range days {
			if log.CreatedAt.Weekday() == d {
				return true
			}
		}
		return false
	}
}

func HourOfDayIn(hours ...int) Condition {
	return func(log models.Log) bool {
		for _, h := range hours {
			if log.CreatedAt.Hour() == h {
				return true
			}
		}
		return false
	}
}

func HourOfDayBetween(min, max int) Condition {
	return func(log models.Log) bool {
		h := log.CreatedAt.Hour()
		return h >= min && h <= max
	}
}

func HourOfDayAfter(min int) Condition {
	return func(log models.Log) bool {
		return log.CreatedAt.Hour() > min
	}
}

func HourOfDayBefore(max int) Condition {
	return func(log models.Log) bool {
		return log.CreatedAt.Hour() < max
	}
}

func IsSeverityIn(severities ...models.Severity) Condition {
	return func(log models.Log) bool {
		for _, s := range severities {
			if log.Severity == s {
				return true
			}
		}
		return false
	}
}

func Not(condition Condition) Condition {
	return func(log models.Log) bool {
		return !condition(log)
	}
}

func And(conditions ...Condition) Condition {
	return func(log models.Log) bool {
		for _, c := range conditions {
			if !c(log) {
				return false
			}
		}
		return true
	}
}

func Or(conditions ...Condition) Condition {
	return func(log models.Log) bool {
		for _, c := range conditions {
			if c(log) {
				return true
			}
		}
		return false
	}
}
