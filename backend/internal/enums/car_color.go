package enums

import (
	"database/sql/driver"
)

type Color string

const (
	RED     Color = "Red"
	BLUE    Color = "Blue"
	YELLOW  Color = "Yellow"
	GREEN   Color = "Green"
	ORANGE  Color = "Orange"
	PURPLE  Color = "Purple"
	BROWN   Color = "Brown"
	BLACK   Color = "Black"
	WHITE   Color = "White"
	GRAY    Color = "Gray"
	CYAN    Color = "Cyan"
	MAGENTA Color = "Magenta"
	LIME    Color = "Lime"
	NAVY    Color = "Navy"
	TEAL    Color = "Teal"
	MAROON  Color = "Maroon"
	OLIVE   Color = "Olive"
	BEIGE   Color = "Beige"
	GOLD    Color = "Gold"
	OTHER   Color = "Other"
)

func (c *Color) Scan(value any) error {
	var sValue string
	switch v := value.(type) {
	case string:
		sValue = v
	case []byte:
		sValue = string(v)
	default:
		return ErrStringConversion
	}
	*c = Color(convertDBFormatToAppFormat(sValue, false))
	return nil
}

func (c Color) Value() (driver.Value, error) {
	return convertAppFormatToDBFormat(string(c)), nil
}

var Colors = []Color{
	RED, BLUE, YELLOW, GREEN, ORANGE, PURPLE, BROWN, BLACK, WHITE, GRAY,
	CYAN, MAGENTA, LIME, NAVY, TEAL, MAROON, OLIVE, BEIGE, GOLD, OTHER}
