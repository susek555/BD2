package scheduler

import "time"

type AuctionEventKind int

const (
	EventAddTimer AuctionEventKind = iota
	EventModifyTimer
	EventForceClose
)

type AuctionEvent struct {
	Kind AuctionEventKind
	At   time.Time
	Cmd  CloseCmd
}
