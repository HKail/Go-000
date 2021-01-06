package main

import (
	"sync"
	"time"
)

// SlidingWindow contains lots of buckets
type SlidingWindow struct {
	mu sync.RWMutex

	winTime     time.Duration
	bucketTime  time.Duration
	bucketCount int
	buckets     []*bucket
}

// bucket content
type bucket struct {
	time  int64
	value float64
}

// NewSlidingWindow return a new sliding window instance
func NewSlidingWindow(winTime time.Duration, bucketTime time.Duration) *SlidingWindow {
	n := new(SlidingWindow)
	n.winTime = winTime
	n.bucketTime = bucketTime
	n.bucketCount = int(n.winTime.Milliseconds() / n.bucketTime.Milliseconds())
	n.buckets = make([]*bucket, 0)

	return n
}

// getCurrentBucket return current bucket of time.Now()
func (sw *SlidingWindow) getCurrentBucket() *bucket {
	now := time.Now().UnixNano() / 1e6

	var b *bucket
	if len(sw.buckets) == 0 {
		b = new(bucket)
		b.time = now
		sw.buckets = append(sw.buckets, b)
		return b
	}

	b = sw.buckets[len(sw.buckets)-1]
	if now-b.time >= int64(sw.bucketTime) {
		b = new(bucket)
		b.time = now
		sw.buckets = append(sw.buckets, b)
	}

	return b
}

// removeOldBuckets remove specific previous bucket
func (sw *SlidingWindow) removeOldBuckets() {
	later := time.Now().UnixNano()/1e6 - int64(sw.winTime)

	for _, b := range sw.buckets {
		if b.time <= later {
			sw.buckets = sw.buckets[1:]
		}
	}
}

// Increment add i to current bucket
func (sw *SlidingWindow) Increment(i float64) {
	if i == 0 {
		return
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	b := sw.getCurrentBucket()
	b.value += i
	sw.removeOldBuckets()
}

// UpdateMax update the maximum value in the current bucket
func (sw *SlidingWindow) UpdateMax(i float64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	b := sw.getCurrentBucket()
	if i > b.value {
		b.value = i
	}
	sw.removeOldBuckets()
}

// Sum sum the values over the buckets in the last specific time
func (sw *SlidingWindow) Sum(now time.Time) float64 {
	sum := float64(0)

	sw.mu.RLock()
	defer sw.mu.RUnlock()

	for _, b := range sw.buckets {
		if b.time >= now.UnixNano()/1e6-int64(sw.winTime) {
			sum += b.value
		}
	}

	return sum
}

// Max return the maximum value seen in last specific time
func (sw *SlidingWindow) Max(now time.Time) float64 {
	var max float64

	sw.mu.RLock()
	defer sw.mu.RUnlock()

	for _, b := range sw.buckets {
		if b.time >= now.UnixNano()/1e6-int64(sw.winTime) {
			if b.value > max {
				max = b.value
			}
		}
	}

	return max
}

// Avg return the average value over the buckets in last specific time
func (sw *SlidingWindow) Avg(now time.Time) float64 {
	return sw.Sum(now) / float64(sw.bucketCount)
}
