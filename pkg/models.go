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
