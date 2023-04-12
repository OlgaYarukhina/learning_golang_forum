package main

import (
	"net/http"
)

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/authentication", app.authenticationHandler)
	mux.HandleFunc("/logout", app.logoutHandler)
	mux.HandleFunc("/my-workspace", app.checkAuth(app.workspaceHandler))
	mux.HandleFunc("/post", app.showPostHandler)
	mux.HandleFunc("/like", app.checkAuth(app.putLike))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
