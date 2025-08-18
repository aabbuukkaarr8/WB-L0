package retry

import (
    "context"
    "math/rand"
    "time"
)

type Options struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    IsRetryable func(error) bool
    OnRetry     func(attempt int, err error)
}

func Do(ctx context.Context, fn func() error, opt Options) error {
    if opt.MaxAttempts <= 0 {
        opt.MaxAttempts = 1
    }
    if opt.BaseDelay <= 0 {
        opt.BaseDelay = 100 * time.Millisecond
    }
    if opt.MaxDelay <= 0 || opt.MaxDelay < opt.BaseDelay {
        opt.MaxDelay = 2 * time.Second
    }

    var delay = opt.BaseDelay
    var lastErr error
    for attempt := 1; attempt <= opt.MaxAttempts; attempt++ {
        if err := fn(); err == nil {
            return nil
        } else {
            lastErr = err
            if attempt == opt.MaxAttempts {
                break
            }
            if opt.IsRetryable != nil && !opt.IsRetryable(err) {
                break
            }
            if opt.OnRetry != nil {
                opt.OnRetry(attempt, err)
            }
            jitter := time.Duration(rand.Int63n(int64(delay)))
            if jitter > opt.MaxDelay {
                jitter = opt.MaxDelay
            }
            select {
            case <-time.After(jitter):
            case <-ctx.Done():
                return ctx.Err()
            }
            if delay < opt.MaxDelay/2 {
                delay *= 2
            }
        }
    }
    return lastErr
}


