package proxy

import (
	"strings"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
)

type ProxyTarget interface {
	Send(log models.Log, rawMessage string) error
}

type ProxyCondition func(log models.Log) bool

type ProxyFilter struct {
	conditions []ProxyCondition
}

type Proxy struct {
	targets ProxyTarget
	filter  ProxyFilter
}

func NewFilter(conditions ...ProxyCondition) ProxyFilter {
	return ProxyFilter{
		conditions: conditions,
	}
}

func (f *ProxyFilter) Evaluate(log models.Log) bool {
	for _, c := range f.conditions {
		if !c(log) {
			return false
		}
	}
	return true
}

func NewProxy(target ProxyTarget, filter ProxyFilter) Proxy {
	return Proxy{
		targets: target,
		filter:  filter,
	}
}

func (p *Proxy) TrySend(log models.Log, rawMessage string) {
	if p.filter.Evaluate(log) {
		p.targets.Send(log, rawMessage)
	}
}

func MessageContains(substr string) ProxyCondition {
	return func(log models.Log) bool {
		return strings.Contains(log.Message, substr)
	}
}

func CategoryContains(substr string) ProxyCondition {
	return func(log models.Log) bool {
		return strings.Contains(log.Category, substr)
	}
}

func ModelContains(substr string) ProxyCondition {
	return func(log models.Log) bool {
		return strings.Contains(log.Model, substr)
	}
}

func IsModel(model string) ProxyCondition {
	return func(log models.Log) bool {
		return log.Model == model
	}
}

func IsCategory(category string) ProxyCondition {
	return func(log models.Log) bool {
		return log.Category == category
	}
}

func IsSeverity(severity models.Severity) ProxyCondition {
	return func(log models.Log) bool {
		return log.Severity == severity
	}
}

func IsSeverityAtLeast(severity models.Severity) ProxyCondition {
	return func(log models.Log) bool {
		return log.Severity >= severity
	}
}

func IsSeverityAtMost(severity models.Severity) ProxyCondition {
	return func(log models.Log) bool {
		return log.Severity <= severity
	}
}

func IsSeverityBetween(min models.Severity, max models.Severity) ProxyCondition {
	return func(log models.Log) bool {
		return log.Severity >= min && log.Severity <= max
	}
}

func TimeBetween(min, max time.Time) ProxyCondition {
	return func(log models.Log) bool {
		return log.CreatedAt.After(min) && log.CreatedAt.Before(max)
	}
}

func TimeAfter(min time.Time) ProxyCondition {
	return func(log models.Log) bool {
		return log.CreatedAt.After(min)
	}
}

func TimeBefore(max time.Time) ProxyCondition {
	return func(log models.Log) bool {
		return log.CreatedAt.Before(max)
	}
}

func DayOfWeekIn(days ...time.Weekday) ProxyCondition {
	return func(log models.Log) bool {
		for _, d := range days {
			if log.CreatedAt.Weekday() == d {
				return true
			}
		}
		return false
	}
}

func HourOfDayIn(hours ...int) ProxyCondition {
	return func(log models.Log) bool {
		for _, h := range hours {
			if log.CreatedAt.Hour() == h {
				return true
			}
		}
		return false
	}
}

func HourOfDayBetween(min, max int) ProxyCondition {
	return func(log models.Log) bool {
		h := log.CreatedAt.Hour()
		return h >= min && h <= max
	}
}

func HourOfDayAfter(min int) ProxyCondition {
	return func(log models.Log) bool {
		return log.CreatedAt.Hour() > min
	}
}

func HourOfDayBefore(max int) ProxyCondition {
	return func(log models.Log) bool {
		return log.CreatedAt.Hour() < max
	}
}

func IsSeverityIn(severities ...models.Severity) ProxyCondition {
	return func(log models.Log) bool {
		for _, s := range severities {
			if log.Severity == s {
				return true
			}
		}
		return false
	}
}

func Not(condition ProxyCondition) ProxyCondition {
	return func(log models.Log) bool {
		return !condition(log)
	}
}

func And(conditions ...ProxyCondition) ProxyCondition {
	return func(log models.Log) bool {
		for _, c := range conditions {
			if !c(log) {
				return false
			}
		}
		return true
	}
}

func Or(conditions ...ProxyCondition) ProxyCondition {
	return func(log models.Log) bool {
		for _, c := range conditions {
			if c(log) {
				return true
			}
		}
		return false
	}
}
