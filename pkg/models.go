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
	Created_at  string
}

type Category struct {
	ID   int
	Name string
}
