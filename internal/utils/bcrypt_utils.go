package utils

type Encryption interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}
