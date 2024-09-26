package format

import (
	"fmt"
	"time"
)

// Returns a string representing a relative date from now (such as "3 days ago")
func FmtRelativeDateToNow(from *time.Time) string {
	now := time.Now()
	if now.Before(*from) {
		return "from the future"
	}

	switch elapsed := now.Unix() - from.Unix(); {
	case elapsed < 60:
		return "just now"
	case elapsed < 60*60:
		plural := "s"
		minutes := elapsed / 60
		if minutes == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d minute%s ago", minutes, plural)
	case elapsed < 60*60*24:
		plural := "s"
		hours := elapsed / (60 * 60)
		if hours == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d hour%s ago", hours, plural)
	case elapsed < 60*60*24*7*2:
		days := elapsed / (60 * 60 * 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	case elapsed < 60*60*24*7*4:
		weeks := elapsed / (60 * 60 * 24 * 7)
		return fmt.Sprintf("%d weeks ago", weeks)
	case elapsed < 60*60*24*365:
		plural := "s"
		months := elapsed / (60 * 60 * 24 * 30)
		if months == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d month%s ago", months, plural)
	default:
		plural := "s"
		years := elapsed / (60 * 60 * 24 * 365)
		if years == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d year%s ago", years, plural)
	}
}
