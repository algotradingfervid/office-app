package clock

import (
	"testing"
	"time"
)

func TestTodayUsesIndiaTime(t *testing.T) {
	// 20:00 UTC on 31 March is already 1 April in India.
	c := Fixed(time.Date(2027, 3, 31, 20, 0, 0, 0, time.UTC))
	if got := Today(c); got != "2027-04-01" {
		t.Fatalf("Today() = %s, want 2027-04-01", got)
	}
}
