package format

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmtRelativeDateToNow(t *testing.T) {
	dates := []time.Time{
		time.Now().Add(60 * time.Second),                      // Now + 1min
		time.Now().Add(-1 * time.Second),                      // Now - 1sec
		time.Now().Add(-10 * time.Second),                     // Now - 10sec
		time.Now().Add(-60 * time.Second),                     // Now - 1min
		time.Now().Add(-10 * 60 * time.Second),                // Now - 10min
		time.Now().Add(-60 * 60 * time.Second),                // Now - 1hr
		time.Now().Add(-10 * 60 * 60 * time.Second),           // Now - 10hr
		time.Now().Add(-24 * 60 * 60 * time.Second),           // Now - 1d
		time.Now().Add(-7 * 24 * 60 * 60 * time.Second),       // Now - 7d
		time.Now().Add(-13 * 24 * 60 * 60 * time.Second),      // Now - 13d
		time.Now().Add(-16 * 24 * 60 * 60 * time.Second),      // Now - 16d
		time.Now().Add(-23 * 24 * 60 * 60 * time.Second),      // Now - 23d
		time.Now().Add(-32 * 24 * 60 * 60 * time.Second),      // Now - 32d
		time.Now().Add(-11 * 30 * 24 * 60 * 60 * time.Second), // Now - 11 months
		time.Now().Add(-23 * 30 * 24 * 60 * 60 * time.Second), // Now - 23 months
		time.Now().Add(-2 * 365 * 24 * 60 * 60 * time.Second), // Now - 2 years
	}
	expected := []string{
		"from the future",
		"just now",
		"just now",
		"1 minute ago",
		"10 minutes ago",
		"1 hour ago",
		"10 hours ago",
		"yesterday",
		"7 days ago",
		"13 days ago",
		"2 weeks ago",
		"3 weeks ago",
		"1 month ago",
		"11 months ago",
		"1 year ago",
		"2 years ago",
	}

	for i, d := range dates {
		assert.Equal(t, expected[i], FmtRelativeDateToNow(&d), fmt.Sprintf("date %v - index %d", d, i))
	}
}
