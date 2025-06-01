package enums

import (
	"database/sql/driver"
)

type Status string

var (
	PENDING   Status = "Pending"
	READY     Status = "Ready"
	PUBLISHED Status = "Published"
	SOLD      Status = "Sold"
)

func (s *Status) Scan(value any) error {
	var sValue string
	switch v := value.(type) {
	case string:
		sValue = v
	case []byte:
		sValue = string(v)
	default:
		return ErrStringConversion
	}
	*s = Status(convertDBFormatToAppFormat(sValue, false))
	return nil
}

func (c Status) Value() (driver.Value, error) {
	return convertAppFormatToDBFormat(string(c)), nil
}
