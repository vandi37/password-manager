package module

type Password struct {
	Id       int64
	Company  string
	Username string
	Password []byte
	Nonce    []byte
	UserId   int64
}
