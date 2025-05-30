package enums

import (
	"database/sql/driver"
	"strings"
)

type Status string

var (
	PENDING   Status = "pending"
	READY     Status = "ready"
	PUBLISHED Status = "published"
)

func (s *Status) Scan(value any) error {
	*s = Status(value.([]byte))
	return nil
}

func (c Status) Value() (driver.Value, error) {
	return strings.ToLower(string(c)), nil
}
