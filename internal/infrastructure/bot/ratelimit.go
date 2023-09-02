package bot

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Limiter interface {
	Limit(ctx context.Context, next func() string) string
}

type Limit struct {
	mu        sync.Mutex
	maxCalls  int
	interval  time.Duration
	lastCalls map[int64][]time.Time
}

func NewRateLimit(ctx context.Context, maxCalls int, interval time.Duration) *Limit {
	return &Limit{
		maxCalls:  maxCalls,
		interval:  interval,
		lastCalls: make(map[int64][]time.Time, maxCalls),
	}
}

// func (l *Limit) Limit(ctx context.Context) {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	if l.maxCalls > len(l.lastCalls) {
// 		l.lastCalls = append(l.lastCalls, time.Now())
// 	}

// 	if time.Since(l.lastCalls[0]) < l.interval {
// 		fmt.Printf("Limit reached, waiting for: %s\n", l.interval-time.Since(l.lastCalls[0]))
// 		time.Sleep(l.interval - time.Since(l.lastCalls[0]))
// 	}

// 	l.lastCalls = l.lastCalls[1:]
// }

func (l *Limit) Limit(ctx context.Context, chatID int64) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if len(l.lastCalls[chatID]) >= l.maxCalls {
		oldestTime := l.lastCalls[chatID][0]
		if now.Sub(oldestTime) < l.interval {
			sleepDuration := l.interval - now.Sub(oldestTime)
			fmt.Println("sleeping for", sleepDuration)
			time.Sleep(sleepDuration)
		}
		l.lastCalls[chatID] = l.lastCalls[chatID][1:]
	}

	l.lastCalls[chatID] = append(l.lastCalls[chatID], now)
}

func (l *Limit) Available(ctx context.Context, chatID int64) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if len(l.lastCalls[chatID]) >= l.maxCalls {
		oldestTime := l.lastCalls[chatID][0]
		if now.Sub(oldestTime) < l.interval {

			return false
		}
		l.lastCalls[chatID] = l.lastCalls[chatID][1:]
	}

	l.lastCalls[chatID] = append(l.lastCalls[chatID], now)

	return true
}

// func (l *Limit) wait(chatID int64) {
// 	now := time.Now()

// 	if len(l.lastCalls[chatID]) >= l.maxCalls {
// 		oldestTime := l.lastCalls[chatID][0]
// 		if now.Sub(oldestTime) < l.interval {
// 			sleepDuration := l.interval - now.Sub(oldestTime)
// 			fmt.Println("sleeping for", sleepDuration)
// 			time.Sleep(sleepDuration)
// 		}
// 		l.lastCalls[chatID] = l.lastCalls[chatID][1:]
// 	}

// 	l.lastCalls[chatID] = append(l.lastCalls[chatID], now)
// }

// func (l *Limit) Limit(ctx context.Context, chatID int64) {
// 	select {
// 	case <-ctx.Done():
// 		return
// 	default:
// 	}

// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	l.wait(chatID)
// }

// func (l *Limit) Available(ctx context.Context, chatID int64) bool {
// 	select {
// 	case <-ctx.Done():
// 		return false
// 	default:
// 	}

// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	l.wait(chatID)

// 	return true
// }
