package main

import (
		"net/http"
		"forum/cmd/web/additional"
	)

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/authentication", app.authentication)
	mux.HandleFunc("/create", app.additional.CreateUser)
	mux.HandleFunc("/show", app.additional.Show)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
