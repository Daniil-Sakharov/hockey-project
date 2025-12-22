package events

import "time"

// PlayerSearchRequested is emitted when user requests player search.
type PlayerSearchRequested struct {
	UserID    int64
	Query     string
	Timestamp time.Time
}

// FilterApplied is emitted when user applies a filter.
type FilterApplied struct {
	UserID     int64
	FilterType string
	Value      string
	Timestamp  time.Time
}

// SessionStarted is emitted when user starts a new session.
type SessionStarted struct {
	UserID    int64
	Timestamp time.Time
}
