package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	ID         int
	username   string
	password   string
	email      string
	created_at time.Time
}

type Post struct {
	ID          int
	Header      string
	Description string
	Category_id int
	User_id     int
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

type Like struct {
	ID      int
	User_id int
	Post_id int
}

type Category struct {
	ID   int
	Name string
}
