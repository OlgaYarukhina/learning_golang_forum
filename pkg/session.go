package models

import "time"

type Session struct {
	Username string
	Expiry   time.Time
}

//проверка действителен ли токен по времени или нет

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
