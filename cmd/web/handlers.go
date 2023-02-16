package main

import (
	"fmt"
	models "forum/pkg"
	"net/http"
	"forum/cmd/additional"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *Application) authentication(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "authent.page.tmpl", &templateData{})
	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
	
		if newUser.additional.ValidateRegistration() == true{
			fmt.Println("All good")

			additional.createUser(newUser, w, r)
			//или сразу в БД но не знаю как педерать павильно модель
			//нужно ли возвращать Id?
			err  := models.Insert(newUser.Username, newUser.Password, newUser.Email)
			if err != nil {

			}
		} else {
			fmt.Println("Something bad")
			// отображение страницы с информацией
			
		}
	}
}

//(m *UserModel) Insert(username, password, email string) (int, error) 

