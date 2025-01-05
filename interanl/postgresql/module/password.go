package module

import "time"

type Password struct {
	Password []byte
	Nonce    []byte
	UserId   int64
	Created  time.Time
}
