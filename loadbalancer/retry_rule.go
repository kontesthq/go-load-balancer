package loadbalancer

import (
	"github.com/hashicorp/consul/api"
	"time"
)

type RetryRule struct {
	subRule        IRule
	maxRetryMillis int64
}

func NewRetryRule(subRule IRule, maxRetryMillis int64) *RetryRule {
	return &RetryRule{
		subRule:        subRule,
		maxRetryMillis: maxRetryMillis,
	}
}

func NewRetryRuleWithDefaults() *RetryRule {
	return NewRetryRule(NewRoundRobinRule(), 500)
}

func (r *RetryRule) SetRule(subRule IRule) {
	r.subRule = subRule
}

func (r *RetryRule) GetRule() IRule {
	return r.subRule
}

func (r *RetryRule) SetMaxRetryMillis(maxRetryMillis int64) {
	r.maxRetryMillis = maxRetryMillis
}

func (r *RetryRule) GetMaxRetryMillis() int64 {
	return r.maxRetryMillis
}

func (r *RetryRule) ChooseServer(lb LoadBalancer) *api.AgentService {
	requestTime := time.Now().UnixMilli()
	deadline := requestTime + r.maxRetryMillis

	var answer *api.AgentService = nil

	answer = r.subRule.ChooseServer(lb)

	for answer == nil && time.Now().UnixMilli() < deadline {
		time.Sleep(1 * time.Millisecond)
		answer = r.subRule.ChooseServer(lb)
	}

	return answer
}
