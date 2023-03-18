package main

import (
	"net/http"
)

func (app *Application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/authentication", app.authentication)
	mux.HandleFunc("/login", app.authorization)
	mux.HandleFunc("/account", app.account)

	mux.HandleFunc("/my-workspace", app.workspace)

	//mux.HandleFunc("/show", app.additional.Show)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
