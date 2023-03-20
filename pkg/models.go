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
	ID         int
	Title      string
	Category   string
	User_id    int
	Content    string
	Created_at time.Time
}
