package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	ID         int
	UserName   string
	Password   string
	Email      string
	Errors     map[string]string
	Created_at time.Time
}

type Post struct {
	ID          int
	Title       string
	Category_id int
	User_id     int
	Content     string
	Like        int
	Created_at  string
}

type Comment struct {
	ID         int
	Comment    string
	Post_id    int
	User_id    int
	Username   string
	Created_at string
}

type Category struct {
	ID   int
	Name string
}
