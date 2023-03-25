package main

import (
	models "forum/pkg"
	"net/http"
)

//проверяем пользователь залогинен или нет, если нет, то мы не даем ему доступ к функции в handler
func (app *Application) checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token") //получаем токен
		if err != nil || err == http.ErrNoCookie {
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		token := c.Value

		userSession, exists := app.Session[token]
		if !exists {
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		if models.Session.IsExpired(userSession) {
			delete(app.Session, token)
			http.Redirect(w, r, "/authentication", 302)
			return
		}
		next(w, r)
	}
}

func (app *Application) yetAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			next(w, r)
			return
		}
		if err == nil || err != http.ErrNoCookie {
			http.Redirect(w, r, "/my-workspace", 302)
			return
		}

		token := c.Value

		userSession, exists := app.Session[token]
		if exists {
			http.Redirect(w, r, "/my-workspace", 302)
			return
		}

		if !models.Session.IsExpired(userSession) {
			http.Redirect(w, r, "/my-workspace", 302)
			return
		}

	}
}
