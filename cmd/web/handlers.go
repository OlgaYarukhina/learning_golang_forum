package main

import (
	"errors"
	"fmt"
	"forum/cmd/web/additional"
	models "forum/pkg"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// array contains session_token + username

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *Application) account(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession := app.Session[sessionToken]

	var objectUser templateData

	users := map[string]string{"username": userSession.Username}
	objectUser.Data = users
	app.render(w, r, ".page.tmpl", &objectUser)
}

func (app *Application) authentication(w http.ResponseWriter, r *http.Request) {

	// save all errors in one variable
	var data templateData

	if r.Method == "POST" {
		newUser := &models.User{
			Email:    r.FormValue("newEmail"),
			UserName: r.FormValue("newUser"),
			Password: r.FormValue("newPassword"),
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := app.Users.CheckUser(email) // get user by email

		if email != "" && password != "" {
			switch additional.CheckPasswordHash(password, user.Password) { //проверяем равен ли пароль который ввел пользователь паролю в БД
			case true:
				//создаем токен
				sessionToken := uuid.NewString()
				//делаем длительность сессии 120 секунд
				expiresAt := time.Now().Add(120 * time.Second)

				//заполняем массив, куда входит токен и имя пользователя
				app.Session[sessionToken] = models.Session{Username: user.UserName, Expiry: expiresAt}

				//устанавливаем куки пользователю и записываем туда имя его и токен
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expiresAt,
				})
			case false:
				data.Data = make(map[string]string)
				data.Data["WrongUserData"] = "Email: " + email + " or Password is wrong! Please, try again"
				fmt.Println("Problem with login")
				app.render(w, r, "authent.page.tmpl", &data)
			}

			if err != nil {
				app.ErrorLog.Println(err)
			}
			http.Redirect(w, r, "/", 303)
			return
		}

		data.Data = additional.ValidateRegistration(newUser)
		fmt.Println(newUser)

		switch len(data.Data) {
		case 0:
			//хешируем пароль
			hashedPassword, err := additional.HashPassword(newUser.Password)
			err = app.Users.Insert(newUser.UserName, hashedPassword, newUser.Email)
			checkUnick := err.Error()
			
			if errors.As(err, &app.sqlError) {
				//app.ErrorLog.Println(err)
				//msg.Data["NewUserExist"] = "User name: " + newUser.UserName + " or Email: " + newUser.Email + " already exist! Please, login or create other user"
				switch checkUnick {
				case "UNIQUE constraint failed: user.email":
					data.Data["NewUserExist"] = "Email " +newUser.Email +" already exists"
				case "UNIQUE constraint failed: user.username":
					data.Data["NewUserExist"] = "User name " + newUser.UserName + " already exists"
				}
				app.render(w, r, "authent.page.tmpl", &data)
			} else {
				// show page with cogratulations or home page with button "Logout"
				app.render(w, r, "home.page.tmpl", &templateData{})
			}

		default:
			app.render(w, r, "authent.page.tmpl", &data)
		}
	}

	if r.Method != "POST" {
		app.render(w, r, "authent.page.tmpl", &data)
	}
}

func (app *Application) workspace(w http.ResponseWriter, r *http.Request) {
	//Получаем токен из куков
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession := app.Session[sessionToken]
	//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен

	//получаем всю информацию из базы данных юзера кому принадлежит этот username

	user, err := app.Users.GetUserByUsername(userSession.Username)
	var data templateData
	data.DataCategories = app.Categories.GetCategories()
	fmt.Println(user.ID)
	data.DataPost = app.Posts.GetUserPosts(user.ID)

	if r.Method == "POST" {
		category := r.FormValue("category")
		categoryId, _ := strconv.Atoi(category)

		newPost := &models.Post{
			User_id:      user.ID,
			Title:        r.FormValue("title"),
			Category_id:  categoryId,
			Content:      r.FormValue("content"),
		}

		err = app.Posts.Insert(newPost.Title, newPost.Content, newPost.Category_id, newPost.User_id)
		if err != nil {
			app.ErrorLog.Println(err)
		}
		http.Redirect(w, r, "/", 303)
		return
	}
	app.render(w, r, "workspace.page.tmpl", &data)
}
