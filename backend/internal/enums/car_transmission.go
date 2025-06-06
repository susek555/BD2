package enums

import (
	"database/sql/driver"
)

type Transmission string

const (
	MANUAL      Transmission = "Manual"
	AUTOMATIC   Transmission = "Automatic"
	CVT         Transmission = "Cvt"
	DUAL_CLUTCH Transmission = "Dual clutch"
)

func (t *Transmission) Scan(value any) error {
	var sValue string
	switch v := value.(type) {
	case string:
		sValue = v
	case []byte:
		sValue = string(v)
	default:
		return ErrStringConversion
	}
	*t = Transmission(convertDBFormatToAppFormat(sValue, false))
	return nil
}

func (t Transmission) Value() (driver.Value, error) {
	return convertAppFormatToDBFormat(string(t)), nil
}

var Transmissions = []Transmission{MANUAL, AUTOMATIC, CVT, DUAL_CLUTCH}
