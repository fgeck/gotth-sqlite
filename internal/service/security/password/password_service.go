package password

import "golang.org/x/crypto/bcrypt"

type PasswordServiceInterface interface {
	HashAndSaltPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type PasswordService struct {
	hashFunc    func(password []byte, cost int) ([]byte, error)
	compareFunc func(hashedPassword, password []byte) error
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		hashFunc:    bcrypt.GenerateFromPassword,
		compareFunc: bcrypt.CompareHashAndPassword,
	}
}
func NewPasswordServiceWithCustomFuncs(
	hashFunc func(password []byte, cost int) ([]byte, error),
	compareFunc func(hashedPassword, password []byte) error,
) *PasswordService {
	return &PasswordService{
		hashFunc:    hashFunc,
		compareFunc: compareFunc,
	}
}

func (s *PasswordService) HashAndSaltPassword(password string) (string, error) {
	hashedPassword, err := s.hashFunc([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *PasswordService) ComparePassword(hashedPassword, password string) error {
	return s.compareFunc([]byte(hashedPassword), []byte(password))
}
