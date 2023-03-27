package main

import (
	"fmt"
	models "forum/pkg"
	"net/http"
)


func (app *Application) checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token") //получаем токен
		if err != nil || err == http.ErrNoCookie {
			data.CheckLogin = false
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		token := c.Value

		userSession, exists := app.Session[token]
		if !exists {
			data.CheckLogin = false
			fmt.Println("Here2")
			fmt.Println(data.CheckLogin)
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		if models.Session.IsExpired(userSession) {
			delete(app.Session, token)
			data.CheckLogin = false
			fmt.Println("Here3")
			fmt.Println(data.CheckLogin)
			http.Redirect(w, r, "/authentication", 302)
			return
		}
		next(w, r)
	}
}



// Вообше ненужная проверка. Если юзер залогирован, у него нет ссылки на авторизацию, как он туда попадет,

// func (app *Application) yetAuth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("session_token")
// 		if err != nil {
// 			next(w, r) // what is this?
// 			return
// 		}
// 		if err == nil || err != http.ErrNoCookie {
// 			http.Redirect(w, r, "/my-workspace", 302)
// 			return
// 		}

// 		token := c.Value

// 		userSession, exists := app.Session[token]
// 		if exists {
// 			http.Redirect(w, r, "/home", 302)
// 			return
// 		}

// 		if !models.Session.IsExpired(userSession) {
// 			http.Redirect(w, r, "/home", 302)
// 			return
// 		}

// 	}
// }
