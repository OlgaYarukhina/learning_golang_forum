package additional

import (
	"fmt"
	models "forum/pkg"
	"regexp"
)


var rxEmail = regexp.MustCompile(`.+@.+\..+`)
var rxUserName = regexp.MustCompile(`.{4,10}`)
var rxPassword = regexp.MustCompile(`.{6,12}`)


func ValidateRegistration(registr *models.User) map[string]string {
	registr.Errors = make(map[string]string)

	

	matchEmail := rxEmail.Match([]byte(registr.Email))
	if matchEmail == false {
		registr.Errors["NewEmail"] = "Please enter a valid email address"
	}

	matchName := rxUserName.Match([]byte(registr.Username))
	if matchName == false {
		registr.Errors["NewUsername"] = "User name must contain at least 4 signs"
	}

	matchPW := rxPassword.Match([]byte(registr.Password))
	if matchPW == false {
		registr.Errors["NewPassword"] = "User password must contain at least 6 signs"
	}
	fmt.Println(registr.Errors)

	return registr.Errors
}
