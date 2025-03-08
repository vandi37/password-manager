package password

import (
	"github.com/vandi37/vanerrors"
	"golang.org/x/crypto/sha3"
)

const (
	ErrorGettingHash = "error getting hash"
)

func Hash(password string, salt []byte) ([]byte, error) {
	hash := sha3.New256()
	_, err := hash.Write(append(salt, []byte(password)...))

	if err != nil {
		return nil, vanerrors.Wrap(ErrorGettingHash, err)
	}

	sum := hash.Sum([]byte{})

	return sum, nil
}
