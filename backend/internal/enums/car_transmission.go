package enums

import (
	"database/sql/driver"
	"strings"
)

type Transmission string

const (
	MANUAL      Transmission = "Manual"
	AUTOMATIC   Transmission = "Automatic"
	CVT         Transmission = "CVT"
	DUAL_CLUTCH Transmission = "Dual clutch"
)

func (t *Transmission) Scan(value any) error {
	*t = Transmission(value.([]byte))
	return nil
}

func (t Transmission) Value() (driver.Value, error) {
	return strings.ToLower(string(t)), nil
}

var Transmissions = []Transmission{MANUAL, AUTOMATIC, CVT, DUAL_CLUTCH}
