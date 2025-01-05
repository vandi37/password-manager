package password

import (
	"fmt"

	"github.com/vandi37/vanerrors"
	"golang.org/x/crypto/sha3"
)

const (
	ErrorGettingHash = "error getting hash"
)

func Hash(password string, salt []byte) (string, error) {
	hash := sha3.New256()
	_, err := hash.Write([]byte(password))

	if err != nil {
		return "", vanerrors.NewWrap(ErrorGettingHash, err, vanerrors.EmptyHandler)
	}

	sha3 := hash.Sum(salt)

	return fmt.Sprintf("%x", sha3), nil
}

func Compare(password, hash string, salt []byte) (bool, error) {
	hashedPassword, err := Hash(password, salt)

	return hash != "" && hashedPassword != "" && hashedPassword == hash, err
}
