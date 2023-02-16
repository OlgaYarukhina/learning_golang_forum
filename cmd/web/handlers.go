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
			}
		} else {
			fmt.Println("Something bad")
			// отображение страницы с информацией

		}
	}
	//на будущее, никогда не ставь app render в самом начале функции
	app.render(w, r, "authent.page.tmpl", &templateData{})
}
