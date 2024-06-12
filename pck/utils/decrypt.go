package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func ValidateUserCredential(savedPassword,inputPassword string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword),[]byte(inputPassword))
	return err == nil
}