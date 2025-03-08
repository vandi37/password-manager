package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/vandi37/vanerrors"
	"golang.org/x/crypto/argon2"
)

const (
	ErrorCreatingCipher = "error creating cipher"
	ErrorCreatingNonce  = "error creating nonce"
	ErrorCreatingGCM    = "error creating GCM"
	ErrorOpeningGCM     = "error opening gcm"
)

func DeriveKey(password, salt []byte) []byte {
	return argon2.Key(password, salt, 3, 32*1024, 4, 32)
}

func Encrypt(data, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, vanerrors.Wrap(ErrorCreatingCipher, err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, vanerrors.Wrap(ErrorCreatingNonce, err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, vanerrors.Wrap(ErrorCreatingGCM, err)
	}
	return aesGCM.Seal(nil, nonce, data, nil), nonce, nil
}

func Decrypt(cipherText, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorCreatingCipher, err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorCreatingGCM, err)
	}
	plaintext, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorOpeningGCM, err)
	}
	return plaintext, nil
}
