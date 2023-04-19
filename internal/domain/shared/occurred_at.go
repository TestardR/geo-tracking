package shared

import (
	"time"
)

const OccurredAtPrecisionLayout = "2006-01-02T15:04:05.999 -0700 MST"

type OccurredAt struct {
	textTime string
}

// OccurredAtFrom will:
// - assure conversion of any TZ to uniform UTC
// - assure equal time precision across all places in application
func OccurredAtFrom(t time.Time) OccurredAt {
	return OccurredAt{
		textTime: t.UTC().Format(OccurredAtPrecisionLayout),
	}
}
