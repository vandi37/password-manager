package password

type PasswordService struct {
	hashSalt  []byte
	argonSalt []byte
}

func New(hashSalt, argonSalt []byte) *PasswordService {
	return &PasswordService{
		hashSalt:  hashSalt,
		argonSalt: argonSalt,
	}
}

func (s *PasswordService) Encrypt(master, data []byte) ([]byte, []byte, error) {
	key := DeriveKey(master, s.argonSalt)
	return Encrypt(data, key)
}

func (s *PasswordService) Decrypt(master, cipherText, nonce []byte) ([]byte, error) {
	key := DeriveKey(master, s.argonSalt)
	return Decrypt(cipherText, key, nonce)
}

func (s *PasswordService) Hash(password string) (string, error) {
	return Hash(password, s.hashSalt)
}

func (s *PasswordService) Compare(password, hash string) (bool, error) {
	return Compare(password, hash, s.hashSalt)
}
