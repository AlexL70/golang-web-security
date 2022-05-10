package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64
}

func (uc *UserClaims) Valid() error {
	if !uc.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token has expired")
	}
	if uc.SessionID == 0 {
		return fmt.Errorf("Invalid session ID")
	}
	return nil
}

var key []byte

func init() {
	for i := byte(1); i <= 64; i++ {
		key = append(key, i)
	}
}

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

func signMessage(msg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("Error signing message: %w", err)
	}
	signature := h.Sum(nil)
	return signature, nil
}

func checkSig(msg, sig []byte) (bool, error) {
	newSig, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("Error checking signature: %w", err)
	}
	same := hmac.Equal(sig, newSig)
	return same, nil
}

func createToken(c *UserClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error in createToken function: %w", err)
	}
	return signedToken, nil
}
