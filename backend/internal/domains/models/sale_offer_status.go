package models

type Status string

var (
	PENDING   Status = "pending"
	READY     Status = "ready"
	PUBLISHED Status = "published"
)
