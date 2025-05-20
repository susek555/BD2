package user

type User struct {
	ID       uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string   `json:"username" gorm:"unique"`
	Password string   `json:"password"`
	Email    string   `json:"email" gorm:"unique"`
	Selector string   `json:"selector" gorm:"check:selector IN ('P', 'C')"`
	Person   *Person  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Company  *Company `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (user *User) GetSubtype() UserSubtype {
	switch user.Selector {
	case "P":
		return user.Person
	case "C":
		return user.Company
	default:
		return nil
	}

}
