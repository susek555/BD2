package scheduler

import "time"

type Item struct {
	AuctionID string
	EndAt     time.Time
	index     int
}

type timerHeap []*Item

func (h timerHeap) Len() int {
	return len(h)
}

func (h timerHeap) Less(i, j int) bool {
	return h[i].EndAt.Before(h[j].EndAt)
}

func (h timerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *timerHeap) Push(x interface{}) {
	*h = append(*h, x.(*Item))
}

func (h *timerHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return item
}
