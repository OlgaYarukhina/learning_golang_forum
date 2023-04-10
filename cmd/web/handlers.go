package main

import (
	"errors"
	"forum/cmd/web/additional"
	models "forum/pkg"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

var msg templateData

//-------------------------- POSTS ------------------------------------------
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.URL.Query().Get("categories") != "" {
		id, err := strconv.Atoi(r.URL.Query().Get("categories"))
		s, err := app.Posts.GetPostsByCategory(id)
		c, err := app.categories.GetCategories()
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "home.page.tmpl", &templateData{
			Posts:      s,
			Categories: c,
		})
	} else {
		s, err := app.Posts.Latest()
		c, err := app.categories.GetCategories()
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "home.page.tmpl", &templateData{
			Posts:      s,
			Categories: c,
		})
	}
}

//  POST -------------------------- Read One post
func (app *Application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	if r.Method == "POST" {
		comment := r.FormValue("comment")
		//Получаем токен из куков
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession := app.Session[sessionToken]
		//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен

		//получаем всю информацию из базы данных юзера кому принадлежит этот username
		user, err := app.Users.GetUserByUsername(userSession.Username)
		if err != nil {
			app.ErrorLog.Println(err)
		}
		app.comment.Insert(comment, id, user.ID, time.Now())
	}
	c, err := app.comment.Comments(id)
	s, err := app.Posts.GetPost(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		Post:    s,
		Comment: c,
	})
}

// POST --------------------------- Create Post
func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	c, err := app.categories.GetCategories()

	if err != nil {
		app.ErrorLog.Println(err)
	}

	if r.Method == "POST" {
		category := r.FormValue("category")
		categoryId, _ := strconv.Atoi(category)
		//Получаем токен из куков
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession := app.Session[sessionToken]
		//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен

		//получаем всю информацию из базы данных юзера кому принадлежит этот username
		user, err := app.Users.GetUserByUsername(userSession.Username)

		newPost := &models.Post{
			Title:      r.FormValue("title"),
			Category:   categoryId,
			User_id:    user.ID,
			Content:    r.FormValue("content"),
			Created_at: time.Now().UTC(),
		}

		err = app.Posts.Insert(newPost.Title, newPost.Content, newPost.Category, newPost.User_id, newPost.Created_at)
		if err != nil {
			app.ErrorLog.Println(err)
		}
		http.Redirect(w, r, "/", 303)
		return
	}

	app.render(w, r, "workspace.page.tmpl", &templateData{
		Categories: c,
	})
}

// SYSTEM OF LIKES
func (app *Application) putLike(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession := app.Session[sessionToken]
	//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен

	//получаем всю информацию из базы данных юзера кому принадлежит этот username
	user, err := app.Users.GetUserByUsername(userSession.Username)

	like := app.Posts.CreateLike(user.ID, id)

	if like != nil {
		app.ErrorLog.Println(like)
	}
	http.Redirect(w, r, "/", 303)
	return
}

// ------------------------------------------- USER ACTION --------------------------------------------

// USER -------- LOGIN + REGISTRATION ------------
func (app *Application) registrationHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	//-----------POST METHOD-------------------
	case "POST":
		newUser := &models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		msg.Data = additional.ValidateRegistration(newUser)
		//-----Create New User
		hashedPassword, err := additional.HashPassword(newUser.Password) // hash password

		err = app.Users.Insert(newUser.Username, hashedPassword, newUser.Email, time.Now())
		//--------------------

		// -------------- IF USERNAME OR EMAIL UNIQUE
		if errors.As(err, &app.sqlError) {
			app.ErrorLog.Println(err)

			if err.Error() == "UNIQUE constraint failed: user.username" {
				msg.Data["Username"] = "User " + newUser.Username + " already exists"
				app.render(w, r, "authent.page.tmpl", &msg)
				return

			}

			if err.Error() == "UNIQUE constraint failed: user.email" {
				msg.Data["Email"] = "User " + newUser.Email + " already exists"
				app.render(w, r, "authent.page.tmpl", &msg)
				return
			}
		}
		// -------------- END IF
		app.render(w, r, "home.page.tmpl", &msg)
		//----------END POST
	case "GET":
		app.render(w, r, "authent.page.tmpl", &msg)
	}
}

//--------- LOGIN ---------------
func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := app.Users.CheckUser(email) // get user by email

		if !additional.CheckPasswordHash(password, user.Password) {
			msg.Data["Errorlogin"] = "Error login"
			app.render(w, r, "authent.page.tmpl", &msg)
			return
		}

		//create Token
		sessionToken := uuid.NewString()
		//120 sec Session
		expiresAt := time.Now().Add(120 * time.Second)
		//array with token and username
		app.Session[sessionToken] = models.Session{Username: user.Username, Expiry: expiresAt}
		//set Cookies
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})

		if err != nil {
			app.ErrorLog.Println(err)
		}
		http.Redirect(w, r, "/", 303)
		return
	//------ END POST
	case "GET":
		app.render(w, r, "authent.page.tmpl", &msg)
	}

}
