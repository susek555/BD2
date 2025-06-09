package scheduler_tests

import (
	"testing"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
)

func TestScheduler_Creation(t *testing.T) {
	scheduler := scheduler.NewScheduler(
		nil, // repo
		nil, // redis
		nil, // notification service
		nil, // sale offer repo
		nil, // purchase creator
		nil, // sale offer service
		nil, // hub
	)

	if scheduler == nil {
		t.Error("Expected scheduler to be created, got nil")
	}
}

func TestScheduler_AddAuction(t *testing.T) {
	scheduler := scheduler.NewScheduler(nil, nil, nil, nil, nil, nil, nil)
	scheduler.AddAuction("123", time.Now().Add(1*time.Hour))
}

func TestScheduler_ForceCloseAuction(t *testing.T) {
	scheduler := scheduler.NewScheduler(nil, nil, nil, nil, nil, nil, nil)
	scheduler.ForceCloseAuction("123", 456, 1000)
}

func TestScheduler_LoadAuctions_NilDependency(t *testing.T) {
	scheduler := scheduler.NewScheduler(nil, nil, nil, nil, nil, nil, nil)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when auctionRetriever is nil, but got none")
		}
	}()

	_ = scheduler.LoadAuctions()
}

func TestScheduler_AddAuctionPastTime(t *testing.T) {
	scheduler := scheduler.NewScheduler(nil, nil, nil, nil, nil, nil, nil)
	scheduler.AddAuction("123", time.Now().Add(-1*time.Hour))
}

func TestScheduler_AddAuctionMultiple(t *testing.T) {
	scheduler := scheduler.NewScheduler(nil, nil, nil, nil, nil, nil, nil)

	scheduler.AddAuction("auction1", time.Now().Add(1*time.Hour))
	scheduler.AddAuction("auction2", time.Now().Add(2*time.Hour))
	scheduler.AddAuction("auction3", time.Now().Add(30*time.Minute))
}
