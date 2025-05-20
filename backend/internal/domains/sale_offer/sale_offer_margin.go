package sale_offer

type MarginValue uint

const (
	LOW_MARGIN    MarginValue = 3
	MEDIUM_MARGIN MarginValue = 5
	HIGH_MARGIN   MarginValue = 10
)

var Margins = []MarginValue{LOW_MARGIN, MEDIUM_MARGIN, HIGH_MARGIN}
