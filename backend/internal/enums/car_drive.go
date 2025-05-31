package enums

import (
	"database/sql/driver"
)

type Drive string

const (
	FWD Drive = "FWD"
	RWD Drive = "RWD"
	AWD Drive = "AWD"
)

func (d *Drive) Scan(value any) error {
	var sValue string
	switch v := value.(type) {
	case string:
		sValue = v
	case []byte:
		sValue = string(v)
	default:
		return ErrStringConversion
	}
	*d = Drive(convertDBFormatToAppFormat(sValue, true))
	return nil
}

func (d Drive) Value() (driver.Value, error) {
	return convertAppFormatToDBFormat(string(d)), nil
}

var Drives = []Drive{FWD, RWD, AWD}
