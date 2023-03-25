package main

import (
	"net/http"
)

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/authentication", app.authentication) 
	mux.HandleFunc("/my-workspace", app.checkAuth(app.workspace))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
