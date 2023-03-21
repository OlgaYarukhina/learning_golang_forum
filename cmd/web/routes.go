package main

import (
	"net/http"
)

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/authentication", app.yetAuth(app.authentication)) //yetAuth - проверяет, если пользователь уже залогинен, то он не пускает его на форму авторизации и логина
	mux.HandleFunc("/account", app.checkAuth(app.account)) //checkAuth - проверяет, залогинен пользователь или нет

	mux.HandleFunc("/my-workspace", app.checkAuth(app.workspace))

	//mux.HandleFunc("/show", app.additional.Show)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
