package user

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
	Selector string `json:"selector" gorm:"check:selector IN ('P', 'C')"`
	Company  *Company
	Person   *Person
}

func (u *User) GetSubtype() UserSubtype {
	switch u.Selector {
	case "P":
		return u.Person
	case "C":
		return u.Company
	default:
		return nil
	}

}
