package enums

type Drive string

const (
	FWD Drive = "FWD"
	RWD Drive = "RWD"
	AWD Drive = "AWD"
)

var Drives = []Drive{FWD, RWD, AWD}
