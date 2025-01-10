package module

type Password struct {
	Id       int
	Company  string
	Username string
	Password []byte
	Nonce    []byte
	UserId   int64
}
