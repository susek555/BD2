package review

type RatingFrequency struct {
	Rating    int `gorm:"column:rating"`
	Frequency int `gorm:"column:frequency"`
}
