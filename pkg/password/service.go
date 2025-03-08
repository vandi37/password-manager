package password

type Service struct {
	hashSalt  []byte
	argonSalt []byte
}

func New(hashSalt, argonSalt []byte) *Service {
	return &Service{
		hashSalt:  hashSalt,
		argonSalt: argonSalt,
	}
}

func (s *Service) Encrypt(master, data []byte) ([]byte, []byte, error) {
	key := DeriveKey(master, s.argonSalt)
	return Encrypt(data, key)
}

func (s *Service) Decrypt(master, cipherText, nonce []byte) ([]byte, error) {
	key := DeriveKey(master, s.argonSalt)
	return Decrypt(cipherText, key, nonce)
}

func (s *Service) Hash(password string) ([]byte, error) {
	return Hash(password, s.hashSalt)
}
