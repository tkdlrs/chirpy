package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	generateHash, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return "", err
	}
	//
	return string(generateHash), nil
}

func CheckPasswordHash(password, hash string) error {
	//
	fmt.Println("CheckPasswordHash ran ")
	//
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	//
	fmt.Println("CheckPasswordHash did not have an error")
	//
	return nil
}
