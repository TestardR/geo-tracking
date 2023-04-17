package system

import (
	"context"
	"time"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

const RetryWaitTimeDefault = time.Second

type Waiter func(sleep time.Duration)
type Retryable func(ctx context.Context) error

type Retry struct {
	retryLimit uint8
	waitTime   time.Duration
	waitCall   Waiter
	logger     shared.StdLogger
}

func NewRetry(retryLimit uint8, waitTime time.Duration, waitCall Waiter, logger shared.StdLogger) *Retry {
	return &Retry{
		retryLimit: retryLimit,
		waitTime:   waitTime,
		logger:     logger,
		waitCall:   waitCall,
	}
}

func (r *Retry) Handle(ctx context.Context, controlled Retryable) error {
	var retry uint8
	var err error
	for err = controlled(ctx); err != nil; {
		r.logger.Printf("Try %d failed with error %s. Waiting for retry", retry, err.Error())
		retry++

		if r.waitTime == RetryWaitTimeDefault {
			r.waitCall(time.Duration(retry) * r.waitTime)
		} else {
			r.waitCall(r.waitTime)
		}
		err = controlled(ctx)
		if retry >= r.retryLimit {
			return err
		}
	}

	return err
}
