package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

var (
	Bcrypt Micro = &micro{}
)

type Micro interface{}

type micro struct{}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
