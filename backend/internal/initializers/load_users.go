package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

func LoadUsers() {
	passwordHash, _ := passwords.Hash("PolskaGurom")

	userCount := 10
	users := make([]models.User, userCount)

	for i := 1; i <= userCount; i++ {
		username := "Herakles" + string(rune('0'+i))
		email := "herakles" + string(rune('0'+i)) + "@gmail.com"
		
		if i%2 == 1 { 
			users[i-1] = models.User{
				Username: username,
				Email:    email,
				Password: passwordHash,
				Selector: "P",
				Person: &models.Person{
					Name:    "Herakles",
					Surname: "Wielki",
				},
			}
		} else { 
			nip := "123456789" + string(rune('0'+(i/2)-1))
			users[i-1] = models.User{
				Username: username,
				Email:    email,
				Password: passwordHash,
				Selector: "C",
				Company: &models.Company{
					Name: username + " Sp. z o.o.",
					Nip:  nip,
				},
			}
		}
	}
	for _, user := range users {
		if err := UserRepo.Create(&user); err != nil {
			panic(err)
		}
	}
}
