package main

import (
	"fmt"
	models "forum/pkg"
	"regexp"
)


var rxEmail = regexp.MustCompile(`.+@.+\..+`)
var rxUserName = regexp.MustCompile(`.{4,10}`)
var rxPassword = regexp.MustCompile(`.{6,12}`)


func  (registr *models.User)validateRegistration() bool {
	registr.Errors = make(map[string]string)
	fmt.Println(registr.Email)
		fmt.Println(registr.Username)
		fmt.Println(registr.Password)

	matchEmail := rxEmail.Match([]byte(registr.Email))
	if matchEmail == false {
		registr.Errors["Email"] = "Please enter a valid email address"
	}

	matchName := rxUserName.Match([]byte(registr.Username))
	if matchName == false {
		registr.Errors["UserName"] = "User name must contain at least 4 signs"
	}

	matchPW := rxPassword.Match([]byte(registr.Password))
	if matchPW == false {
		registr.Errors["Password"] = "User password contain at least 6 signs"
	}
	fmt.Println(registr.Errors)

	return len(registr.Errors) == 0
}


func (registr *models.User) validateLogin() bool {
	registr.Errors = make(map[string]string)



	return len(registr.Errors) == 0
}