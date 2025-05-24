package models

type Status string

const (
	PENDING   Status = "pending"
	READY     Status = "ready"
	PUBLISHED Status = "published"
)
