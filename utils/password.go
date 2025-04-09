package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword gera um hash da senha
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
