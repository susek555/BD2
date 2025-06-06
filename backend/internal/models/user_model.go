package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Selector string `json:"selector" gorm:"type:SELECTOR"`
	Person   *Person
	Company  *Company
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
