package proxy

import (
	"sadk.dev/logar/logfilter"
	"sadk.dev/logar/models"
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
