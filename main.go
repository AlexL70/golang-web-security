package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pass := "123456789"
	hashedPass, err := hashPassword(pass)
	if err != nil {
		panic(err)
	}
	//pass = "12345678"
	err = comparePassword(pass, hashedPass)
	if err != nil {
		log.Fatalln("Not logged in.", err)
	}
	log.Println("Logged in!")
}

func hashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting password: %w", err)
	}
	return bs, nil
}

func comparePassword(passowrd string, hashedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(passowrd))
	if err != nil {
		return fmt.Errorf("Invalid password: %w", err)
	}
	return nil
}
