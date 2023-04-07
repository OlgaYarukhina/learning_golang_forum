package main

import (
	"errors"
	models "forum/pkg"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.URL.Query().Get("categories") != "" {
		id, err := strconv.Atoi(r.URL.Query().Get("categories"))
		s, err := app.posts.GetPostsByCategory(id)

		c, err := app.categories.GetCategories()

		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "home.page.tmpl", &templateData{
			Posts:      s,
			Categories: c,
		})
	} else {
		s, err := app.posts.Latest()

		c, err := app.categories.GetCategories()

		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "home.page.tmpl", &templateData{
			Posts:      s,
			Categories: c,
		})
	}

}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	c, err := app.comment.Comments(id)
	s, err := app.posts.GetPost(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Post:    s,
		Comment: c,
	})
}
