package schema

import "time"

type HistoryRecord struct {
	UserId    int
	Segment   string
	Action    string
	Timestamp time.Time
}
