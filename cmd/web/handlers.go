package main

import (
	"errors"
	"fmt"
	"forum/cmd/web/additional"
	models "forum/pkg"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var sessions = map[string]models.Session{} // array contains session_token + username

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *Application) account(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")

	if err != nil || err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", 302)
		return
	}

	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]

	if !exists {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if models.Session.IsExpired(userSession) {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", 302)
		return
	}
	var objectUser templateData

	users := map[string]string{"username": userSession.Username}
	objectUser.Data = users
	app.render(w, r, "account.page.tmpl", &objectUser)
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

		msg.Data = additional.ValidateRegistration(newUser)

		switch len(msg.Data) {
		case 0:
			// show page with cogratulations or home page with button "Logout"
			app.render(w, r, "home.page.tmpl", &templateData{})
			hashedPassword, err := additional.HashPassword(newUser.Password)

			//модель хранится в app, если ты работаешь с моделями, то только в handler работай
			err = app.Users.Insert(newUser.Username, hashedPassword, newUser.Email)
			if errors.As(err, &app.sqlError) {
				app.ErrorLog.Println(err)
				// return wich fild is not unic
				// add check errors

				switch errors.Is(err, errors.New("UNIQUE constraint failed: user.username")) {
				case true:
					msg.Data["Username"] = "User " + newUser.Username + " already exists"
				case false:
					msg.Data["Email"] = "Email " + newUser.Email + " already exists"
				}

				app.render(w, r, "authent.page.tmpl", &msg)
			}
		default:
			app.render(w, r, "authent.page.tmpl", &msg)
		}
	} else {
		app.render(w, r, "authent.page.tmpl", &msg)
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
		//на будущее, никогда не ставь app render в самом начале функции
		app.render(w, r, "authent.page.tmpl", &templateData{})
	}

	app.render(w, r, "workspace.page.tmpl", &templateData{})
}

func (app *Application) authorization(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		email := r.FormValue("email")
		fmt.Println(email)
		password := r.FormValue("password")
		user, err := app.Users.CheckUser(email) // get user by email

		switch additional.CheckPasswordHash(password, user.Password) {
		case true:
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(120 * time.Second)

			sessions[sessionToken] = models.Session{Username: user.Username, Expiry: expiresAt}

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})
		case false:
			fmt.Println("bad")
		}

		if err != nil {
			app.ErrorLog.Println(err)
		}
		app.render(w, r, "authent.page.tmpl", &templateData{})
	}
	if r.Method != "POST" {
		//на будущее, никогда не ставь app render в самом начале функции
		app.render(w, r, "authent.page.tmpl", &templateData{})
	}
}
