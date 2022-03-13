package libtime

import (
	"time"
)

// StopFunc is used to stop a time.Timer.
//
// Calling StopFunc prevents its time.Timer from firing. Returns true if the call
// stops the timer, false if the timer has already expired. or has been stopped.
//
// https://pkg.go.dev/time#Timer.Stop
type StopFunc func() bool

// SafeTimer creates a time.Timer and a StopFunc, forcing the caller to deal
// with the otherwise potential resource leak. Encourages safe use of a time.Timer
// in a select statement, but without the overhead of a context.Context.
//
// Typical usage:
//
//    t, stop := libtime.SafeTimer(interval)
//    defer stop()
//    for {
//      select {
//        case <- t.C:
//          foo()
//        case <- otherC :
//          return
//      }
//    }
//
// Does not panic if duration is <= 0, instead assuming the smallest positive value.
func SafeTimer(duration time.Duration) (*time.Timer, StopFunc) {
	t := time.NewTimer(duration)
	return t, t.Stop
}
