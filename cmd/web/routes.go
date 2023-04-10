package main

import "net/http"

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/registration", app.yetAuth(app.registrationHandler)) //yetAuth - проверяет, если пользователь уже залогинен, то он не пускает его на форму авторизации и логина
	mux.HandleFunc("/login", app.yetAuth(app.loginHandler))               //checkAuth - проверяет, залогинен пользователь или нет

	mux.HandleFunc("/createPost", app.checkAuth(app.createPostHandler))

	//mux.HandleFunc("/show", app.additional.Show)
	mux.HandleFunc("/post", app.showPostHandler)

	mux.HandleFunc("/like", app.checkAuth(app.putLike))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
