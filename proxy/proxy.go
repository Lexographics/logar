package proxy

import (
	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/logfilter"
)

type ProxyTarget interface {
	Send(log models.Log, rawMessage string) error
}

type Proxy struct {
	target ProxyTarget
	filter logfilter.Filter
}

func NewProxy(target ProxyTarget, filter logfilter.Filter) Proxy {
	return Proxy{
		target: target,
		filter: filter,
	}
}

func (p *Proxy) TrySend(log models.Log, rawMessage string) {
	if p.filter.Evaluate(log) {
		p.target.Send(log, rawMessage)
	}
}
