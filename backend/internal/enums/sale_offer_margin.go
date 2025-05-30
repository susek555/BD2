package enums

import "database/sql/driver"

type MarginValue uint

const (
	LOW_MARGIN    MarginValue = 3
	MEDIUM_MARGIN MarginValue = 5
	HIGH_MARGIN   MarginValue = 10
)

func (m *MarginValue) Scan(value any) error {
	*m = MarginValue(value.(int64))
	return nil
}

func (m MarginValue) Value() (driver.Value, error) {
	return int64(m), nil
}

var Margins = []MarginValue{LOW_MARGIN, MEDIUM_MARGIN, HIGH_MARGIN}
