package main

import (
	"fmt"
	"time"
)

// JobTicker ticks at every 10pm
type JobTicker struct {
	t *time.Timer
}

const (
	// HourToTick stores hour value for ticking time
	HourToTick int = 22
	// MinuteToTick stores minute value for ticking time
	MinuteToTick int = 00
	// SecondToTick stores second value for ticking time
	SecondToTick int = 00
	// IntervalPeriod stores duration for a day
	IntervalPeriod time.Duration = 24 * time.Hour
)

// NewJobTicker creates a job ticker instance
func NewJobTicker() JobTicker {
	return JobTicker{time.NewTimer(getNextTickDuration())}
}

func getNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), HourToTick, MinuteToTick, SecondToTick, 0, time.Local)

	if nextTick.Before(time.Now()) {
		nextTick = nextTick.Add(IntervalPeriod)
	}

	fmt.Printf("next tick %v\n", nextTick.String())
	return nextTick.Sub(time.Now())
}

func (jt *JobTicker) updateJobTicker() {
	jt.t.Reset(getNextTickDuration())
}
