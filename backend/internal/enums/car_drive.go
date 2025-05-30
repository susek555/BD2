package enums

import (
	"database/sql/driver"
	"strings"
)

type Drive string

const (
	FWD Drive = "FWD"
	RWD Drive = "RWD"
	AWD Drive = "AWD"
)

func (d *Drive) Scan(value any) error {
	*d = Drive(value.([]byte))
	return nil
}

func (d Drive) Value() (driver.Value, error) {
	return strings.ToLower(string(d)), nil
}

var Drives = []Drive{FWD, RWD, AWD}
