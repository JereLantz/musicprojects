package utils

import "time"

type Credentials struct{
	Username string
	Password string
}

type Session struct {
	LoggedIn bool
	Username string
	Expiry time.Time
}


func (s Session) IsSessionExpired() bool{
	return s.Expiry.Before(time.Now())
}

