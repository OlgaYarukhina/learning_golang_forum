package main

import (
		"net/http"
		"forum/cmd/web/additional"
	)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/authentication", app.authentication)
	mux.HandleFunc("/create", app.additional.create)
	mux.HandleFunc("/show", app.additional.show)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
