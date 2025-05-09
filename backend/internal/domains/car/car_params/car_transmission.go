package car_params

type Transmission string

const (
	MANUAL      Transmission = "Manual"
	AUTOMATIC   Transmission = "Automatic"
	CVT         Transmission = "CVT"
	DUAL_CLUTCH Transmission = "Dual clutch"
)

var Transmissions = []Transmission{MANUAL, AUTOMATIC, CVT, DUAL_CLUTCH}
