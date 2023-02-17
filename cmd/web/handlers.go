package main

import (
	"fmt"
	"forum/cmd/web/additional"
	models "forum/pkg"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *Application) authentication(w http.ResponseWriter, r *http.Request) {

	// save all errors in one variable
	var msg *templateData


	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		if additional.ValidateRegistration(newUser) == true {
			fmt.Println("All good")
			
			//модель хранится в app, если ты работаешь с моделями, то только в handler работай
			err := app.Users.Insert(newUser.Username, newUser.Email, newUser.Password)
			if err != nil {
				app.ErrorLog.Println()
				// return wich fild is not unic
				// add check errors
				msg.Errors["Username"]  = "User already exists"
			    msg.Errors["Email"]  = "Email already exists"
				app.render(w, r, "authent.page.tmpl", msg)

			}
			// show page with cogratulations or home page with button "Logout"
			
		} else {
			fmt.Println("Something bad")
			// show page with wrong data registrations
			
			// problems here
			msg.Errors["Username"] = newUser.Username
			msg.Errors["Email"] = newUser.Email
			
			fmt.Print("Here")
			app.render(w, r, "authent.page.tmpl", msg)

		}
	}

	//на будущее, никогда не ставь app render в самом начале функции
	app.render(w, r, "authent.page.tmpl", &templateData{})
}
