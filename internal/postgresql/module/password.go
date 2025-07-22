package module

type Password struct {
	Id       int
	Company  string
	Username string
	Password []byte
	Nonce    []byte
	UserId   int64
}

type User struct {
	Id       int64
	Password []byte
	Key      []byte
	Nonce    []byte
}
