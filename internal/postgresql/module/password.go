package module

import "time"

type Password struct {
	Company  string
	Username string
	Password []byte
	Nonce    []byte
	UserId   int64
	Created  time.Time
}
