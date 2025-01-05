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
		return nil, nil, vanerrors.NewWrap(ErrorCreatingCipher, err, vanerrors.EmptyHandler)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, vanerrors.NewWrap(ErrorCreatingNonce, err, vanerrors.EmptyHandler)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, vanerrors.NewWrap(ErrorCreatingGCM, err, vanerrors.EmptyHandler)
	}
	return aesGCM.Seal(nil, nonce, data, nil), nonce, nil
}

func Decrypt(cipherText, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingCipher, err, vanerrors.EmptyHandler)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingGCM, err, vanerrors.EmptyHandler)
	}
	plaintext, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningGCM, err, vanerrors.EmptyHandler)
	}
	return plaintext, nil
}
