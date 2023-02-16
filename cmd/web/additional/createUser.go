package additional

import (
	"fmt"
	models "forum/pkg"
	"net/http"
	"strconv"
)

func (app *Application) CreateUser(user *models.User, w http.ResponseWriter, r *http.Request) {
	name := user.Username
	email := user.Email
	password := user.Password

	id, err := app.users.Insert(name, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Println(id)
	w.Write([]byte(string(id)))
}

func (app *Application) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Отображение ID %d...", id)
}
