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
	var msg *templateData
	
	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		checkValid := additional.ValidateRegistration(newUser)

		if len(checkValid) == 0 {
			// show page with cogratulations or home page with button "Logout"
			app.render(w, r, "home.page.tmpl", &templateData{})
			
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
			
		} else {
			for key, value := range checkValid {
				
				if key == "Username" {
					msg.Errors["Username"] = value
				}
				if key == "Email" {
					msg.Errors["Email"] = value
				}
				if key == "Password" {
					msg.Errors["Password"] = value
				}
			}
			app.render(w, r, "authent.page.tmpl", msg)

		}
	}

	if r.Method != "POST" {
			//на будущее, никогда не ставь app render в самом начале функции
	app.render(w, r, "authent.page.tmpl", &templateData{})
		}


}
