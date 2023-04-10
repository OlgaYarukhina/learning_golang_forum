package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	ID         int
	Username   string
	Password   string
	Email      string
	Errors     map[string]string
	Created_at time.Time
}

type Post struct {
	ID         int
	Title      string
	Category   int
	User_id    int
	Like       int
	Content    string
	Created_at time.Time
}

type Comment struct {
	ID         int
	Comment    string
	Post_id    int
	User_id    int
	Username   string
	Created_at string
}

type Like struct {
	ID      int
	User_id int
	Post_id int
}

type Category struct {
	ID   int
	Name string
}
