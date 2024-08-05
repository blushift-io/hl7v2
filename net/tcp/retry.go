package tcp

import "time"

type Retry interface {
	Backoff(uint64) (time.Duration, bool)
}

type expBackoffRetry struct {
	Count        int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func (r expBackoffRetry) Backoff(count uint64) (time.Duration, bool) {
	if count > uint64(r.Count) {
		return r.MaxDelay, false
	}

	delay := r.InitialDelay * (1 << count)
	if delay > r.MaxDelay {
		delay = r.MaxDelay
	}

	return delay, true
}
