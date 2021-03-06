package ratelimiter

import (
	kcache "k8s.io/kubernetes/pkg/client/cache"
	kutil "k8s.io/kubernetes/pkg/util"
	utilruntime "k8s.io/kubernetes/pkg/util/runtime"
	utilwait "k8s.io/kubernetes/pkg/util/wait"
)

// HandlerFunc defines function signature for a RateLimitedFunction.
type HandlerFunc func() error

// RateLimitedFunction is a rate limited function controlling how often the function/handler is invoked.
type RateLimitedFunction struct {
	// Handler is the function to rate limit calls to.
	Handler HandlerFunc

	// Internal queue of requests to be processed.
	queue kcache.Queue

	// Rate limiting configuration.
	kutil.RateLimiter
}

// NewRateLimitedFunction creates a new rate limited function.
func NewRateLimitedFunction(keyFunc kcache.KeyFunc, interval int, handlerFunc HandlerFunc) *RateLimitedFunction {
	fifo := kcache.NewFIFO(keyFunc)

	qps := float32(1000.0) // Call rate per second (SLA).
	if interval > 0 {
		qps = float32(1.0 / float32(interval))
	}

	limiter := kutil.NewTokenBucketRateLimiter(qps, 1)

	return &RateLimitedFunction{handlerFunc, fifo, limiter}
}

// RunUntil begins processes the resources from queue asynchronously until
// stopCh is closed.
func (rlf *RateLimitedFunction) RunUntil(stopCh <-chan struct{}) {
	go utilwait.Until(func() { rlf.handleOne(rlf.queue.Pop()) }, 0, stopCh)
}

// handleOne processes a request in the queue invoking the rate limited
// function.
func (rlf *RateLimitedFunction) handleOne(resource interface{}) {
	rlf.RateLimiter.Accept()
	if err := rlf.Handler(); err != nil {
		utilruntime.HandleError(err)
	}
}

// Invoke adds a request if its not already present and waits for the
// background processor to execute it.
func (rlf *RateLimitedFunction) Invoke(resource interface{}) {
	rlf.queue.AddIfNotPresent(resource)
}
