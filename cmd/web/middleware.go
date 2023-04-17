package main

import (
	"fmt"
	models "forum/pkg"
	"net/http"
	"regexp"
	"strings"
)

var browserName = ""

func (app *Application) checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if browserName == "" {
			browserName, _ = getBrowserName(r)
		} else {
			browserNameCurrent, _ := getBrowserName(r)
			if browserNameCurrent != browserName {
				http.Redirect(w, r, "/", 302)
				return
			}
		}

		c, err := r.Cookie("session_token") //получаем токен
		if err != nil || err == http.ErrNoCookie {
			data.CheckLogin = false
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		token := c.Value

		userSession, exists := app.Session[token]
		if !exists {
			data.CheckLogin = false
			fmt.Println("Here2")
			fmt.Println(data.CheckLogin)
			http.Redirect(w, r, "/authentication", 302)
			return
		}

		if models.Session.IsExpired(userSession) {
			delete(app.Session, token)
			data.CheckLogin = false
			fmt.Println("Here3")
			fmt.Println(data.CheckLogin)
			http.Redirect(w, r, "/authentication", 302)
			return
		}
		next(w, r)
	}
}

func getBrowserName(r *http.Request) (name, version string) {
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Trident/") {
		name = "Internet Explorer"
		if strings.Contains(userAgent, "MSIE") {
			re := regexp.MustCompile("MSIE ([0-9]+\\.[0-9]+)")
			matches := re.FindStringSubmatch(userAgent)
			if len(matches) > 1 {
				version = matches[1]
			}
		} else {
			re := regexp.MustCompile("rv:([0-9]+\\.[0-9]+)")
			matches := re.FindStringSubmatch(userAgent)
			if len(matches) > 1 {
				version = matches[1]
			}
		}
	} else if strings.Contains(userAgent, "Firefox") {
		name = "Firefox"
		re := regexp.MustCompile("Firefox/([0-9]+\\.[0-9]+)")
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			version = matches[1]
		}
	} else if strings.Contains(userAgent, "Chrome") {
		name = "Chrome"
		re := regexp.MustCompile("Chrome/([0-9]+\\.[0-9]+)")
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			version = matches[1]
		}
	} else if strings.Contains(userAgent, "Safari") {
		name = "Safari"
		re := regexp.MustCompile("Version/([0-9]+\\.[0-9]+)")
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			version = matches[1]
		}
	} else {
		name = "Unknown"
	}
	return
}

// Вообше ненужная проверка. Если юзер залогирован, у него нет ссылки на авторизацию, как он туда попадет,

// func (app *Application) yetAuth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("session_token")
// 		if err != nil {
// 			next(w, r) // what is this?
// 			return
// 		}
// 		if err == nil || err != http.ErrNoCookie {
// 			http.Redirect(w, r, "/my-workspace", 302)
// 			return
// 		}

// 		token := c.Value

// 		userSession, exists := app.Session[token]
// 		if exists {
// 			http.Redirect(w, r, "/home", 302)
// 			return
// 		}

// 		if !models.Session.IsExpired(userSession) {
// 			http.Redirect(w, r, "/home", 302)
// 			return
// 		}

// 	}
// }
