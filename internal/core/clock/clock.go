// Package clock gives every date calculation one notion of "now", in India time.
// Tests replace it with Fixed so "today" is known.
package clock

import "time"

// IST is Asia/Kolkata. Set in code, never read from the machine's TZ.
var IST = mustLoad("Asia/Kolkata")

// Clock returns the current time in IST.
type Clock interface {
	Now() time.Time
}

// System is the real clock.
type System struct{}

// Now returns the current time in IST.
func (System) Now() time.Time { return time.Now().In(IST) }

// Fixed always returns the same instant, in IST.
type Fixed time.Time

// Now returns the fixed instant in IST.
func (f Fixed) Now() time.Time { return time.Time(f).In(IST) }

// Today returns c's current calendar date as YYYY-MM-DD.
func Today(c Clock) string { return c.Now().Format(time.DateOnly) }

func mustLoad(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
