package enums

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

var Colors = []Color{
	RED, BLUE, YELLOW, GREEN, ORANGE, PURPLE, BROWN, BLACK, WHITE, GRAY,
	CYAN, MAGENTA, LIME, NAVY, TEAL, MAROON, OLIVE, BEIGE, GOLD, OTHER}
