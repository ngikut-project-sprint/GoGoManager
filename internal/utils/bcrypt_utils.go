package utils

type HashPassword func(password []byte, cost int) ([]byte, error)

type ComparePassword func(hashedPassword, password []byte) error

type Encryption interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}
