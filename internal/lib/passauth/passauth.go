package passauth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BcryptProvider interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type RealBcrypt struct{}

func (rb RealBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (rb RealBcrypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

type PassAuthService struct {
	bcrypt      BcryptProvider
	defaultCost int
}

func ExecAuthService(rb BcryptProvider) *PassAuthService {
	return &PassAuthService{
		bcrypt:      rb,
		defaultCost: 14,
	}
}

// Генерация пароля с хешем в БД
func (pas *PassAuthService) GeneratePassword(password []byte) ([]byte, error) {
	// Сложность хеширования пароля
	hash, err := pas.bcrypt.GenerateFromPassword(password, pas.defaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate pass %v", err)
	}
	return hash, nil
}

// Проверка хеш пароля на валидность
func (pas *PassAuthService) VerifyPassword(password, passHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return status.Errorf(codes.InvalidArgument, "invalid password")
		}
		return status.Errorf(codes.Internal, "failed to verify password %v", err)
	}
	return nil
}
