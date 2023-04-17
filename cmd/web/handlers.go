package main

import (
	"errors"
	"forum/cmd/web/additional"
	models "forum/pkg"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var data templateData

// array contains session_token + username

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.URL.Query().Get("categories") != "" {
		id, err := strconv.Atoi(r.URL.Query().Get("categories"))
		s, err := app.Posts.GetPostsByCategory(id)
		c := app.Categories.GetCategories()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data.DataCategories = c
		data.DataPost = s
		app.render(w, r, "home.page.tmpl", &data)
	} else {
		s, err := app.Posts.Latest()
		c := app.Categories.GetCategories()
		data.DataCategories = c
		data.DataPost = s
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "home.page.tmpl", &data)
	}

}

func (app *Application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	if r.Method == "POST" {
		c, err := r.Cookie("session_token") //получаем токен
		if err != nil || err == http.ErrNoCookie {
			data.CheckLogin = false
			http.Redirect(w, r, "/authentication", 302)
			return
		}
		comment := r.FormValue("comment")
		//Получаем токен из куков
		sessionToken := c.Value
		userSession := app.Session[sessionToken]
		//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен
		//получаем всю информацию из базы данных юзера кому принадлежит этот username
		user, err := app.Users.GetUserByUsername(userSession.Username)
		if err != nil {
			app.ErrorLog.Println(err)
		}
		app.Comment.Insert(comment, id, user.ID, time.Now())
	}
	c, err := app.Comment.Comments(id)
	s, err := app.Posts.GetPost(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data.DataComment = c
	data.SinglePost = s
	app.render(w, r, "show.page.tmpl", &data)
}

func (app *Application) authenticationHandler(w http.ResponseWriter, r *http.Request) {

	data.Data = make(map[string]string)

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email != "" || password != "" {
			user, err := app.Users.CheckUser(email)                        // get user by email
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
				data.Data["WrongUserData"] = "Email: " + email + " or Password is wrong! Please, try again"
				app.render(w, r, "authent.page.tmpl", &data)
			}

			if err != nil {
				app.ErrorLog.Println(err)
			}
			data.CheckLogin = true
			http.Redirect(w, r, "/", 303)
			return
		} else {

			newUser := &models.User{
				Email:    r.FormValue("newEmail"),
				UserName: r.FormValue("newUser"),
				Password: r.FormValue("newPassword"),
			}

			data.Data = additional.ValidateRegistration(newUser)

			if len(data.Data) == 0 {
				//хешируем пароль
				hashedPassword, errHash := additional.HashPassword(newUser.Password)
				if errHash != nil {
					app.ErrorLog.Println(errHash)
				}
				err := app.Users.Insert(newUser.UserName, hashedPassword, newUser.Email)
				if errors.As(err, &app.sqlError) {
					checkUnick := err.Error()
					switch checkUnick {
					case "UNIQUE constraint failed: user.email":
						data.Data["NewUserExist"] = "Email " + newUser.Email + " already exists"
					case "UNIQUE constraint failed: user.username":
						data.Data["NewUserExist"] = "User name " + newUser.UserName + " already exists"
					}
					app.render(w, r, "authent.page.tmpl", &data)
				} else {
					// show page with cogratulations or home page with button "Logout"
					data.Data["UserWasCreate"] = "User " + newUser.UserName + " created. Please login"
					app.render(w, r, "authent.page.tmpl", &data)
				}
			} else {
				app.render(w, r, "authent.page.tmpl", &data)
			}
		}
	} else {
		app.render(w, r, "authent.page.tmpl", &data)
	}
}

func (app *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	delete(app.Session, "UserID")
	data.CheckLogin = false
	http.Redirect(w, r, "/authentication", 302)

}

func (app *Application) workspaceHandler(w http.ResponseWriter, r *http.Request) {
	//Получаем токен из куков
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession := app.Session[sessionToken]
	//userSeesion - хранит в себе userName конкретного пользователю кому принадлежит сам токен
	//получаем всю информацию из базы данных юзера кому принадлежит этот username
	user, err := app.Users.GetUserByUsername(userSession.Username)

	data.Data = make(map[string]string)
	data.DataCategories = app.Categories.GetCategories()
	data.DataPost = app.Posts.GetUserPosts(user.ID)
	data.DataPostAdditional = app.Posts.GetUserFavoritePosts(user.ID)

	if r.Method == "POST" {
		r.ParseForm()
		category := r.Form["category_type"]
		newPost := &models.Post{
			User_id:       user.ID,
			Title:         r.FormValue("title"),
			Category_name: category,
			Content:       r.FormValue("content"),
		}

		err = app.Posts.Insert(newPost.Title, newPost.Content, newPost.User_id, time.Now(), category)
		if err != nil {
			data.Data["PostWasCreated"] = "Post was not created!"
			app.render(w, r, "workspace.page.tmpl", &data)
			app.ErrorLog.Println(err)
		}
		data.Data["PostWasCreated"] = "Post was created! Please, refresh page to see post"
		app.render(w, r, "workspace.page.tmpl", &data)

		return
	}
	data.Data["PostWasCreated"] = ""
	app.render(w, r, "workspace.page.tmpl", &data)
}

// SYSTEM OF LIKES
func (app *Application) putLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("comment_id") != "" {
		id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
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
		err = app.Comment.CreateLike(id, user.ID)
		if err != nil {
			app.ErrorLog.Println(err)
		}
		http.Redirect(w, r, r.Header.Get("Referer"), 303)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	type_of_like := r.URL.Query().Get("type")
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
	err = app.Posts.CreateLike(user.ID, id, type_of_like)
	if err != nil {
		app.ErrorLog.Println(err)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 303)
	return
}

// <p>Hohd shift to add 2 or more categories</p>
// <select class="form-select" id="category1" name="category1" multiple aria-label="multiple select example" required>
//   <option selected>Choose category</option>
//   {{range .DataCategories}}
//   <option value="{{.ID}}">{{.Name}}</option>
//   {{end}}
// </select>
