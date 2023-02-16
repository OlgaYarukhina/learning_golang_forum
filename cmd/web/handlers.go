package main

import (
	"fmt"
	models "forum/pkg"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *application) authentication(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "authent.page.tmpl", &templateData{})
	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
	
		if newUser.validateRegistration() == true{
			fmt.Println("All good")
		} else {
			fmt.Println("Something bad")
		}
	}
}



func (app *application) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	name := "Christian"
	email := "email@mail.ru"
	password := "123123"

	id, err := app.users.Insert(name, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Println(id)
	w.Write([]byte(string(id)))
}

func (app *application) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Отображение ID %d...", id)
}
