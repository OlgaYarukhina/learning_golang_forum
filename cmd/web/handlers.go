package main

import (
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
	var msg templateData

	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		msg.Errors = additional.ValidateRegistration(newUser)

		if len(msg.Errors) == 0 {
			// show page with cogratulations or home page with button "Logout"
			app.render(w, r, "home.page.tmpl", &templateData{})

			//модель хранится в app, если ты работаешь с моделями, то только в handler работай
			err := app.Users.Insert(newUser.Username, newUser.Email, newUser.Password)
			if err != nil {
				app.ErrorLog.Println()
				// return wich fild is not unic
				// add check errors
				msg.Errors["Username"] = "User " + newUser.Username + " already exists"
				msg.Errors["Email"] = "Email " + newUser.Email + " already exists"
				app.render(w, r, "authent.page.tmpl", &msg)
			}
		} else {
			app.render(w, r, "authent.page.tmpl", &msg)
		}
	}
	if r.Method != "POST" {
		app.render(w, r, "authent.page.tmpl", &templateData{})
	}
}

func (app *Application) workspace(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		newPost := &models.Post{
			Title:    r.FormValue("title"),
			Category: r.FormValue("category"),
			Content:  r.FormValue("content"),
		}

		app.render(w, r, "workspace.page.tmpl", &templateData{})
		err := app.Posts.Insert(newPost.Title, newPost.Content)
		if err != nil {
			app.ErrorLog.Println()
		}
	}

	if r.Method != "POST" {
		app.render(w, r, "workspace.page.tmpl", &templateData{})
	}
}
