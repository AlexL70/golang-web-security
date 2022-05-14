package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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

type key struct {
	key     []byte
	created time.Time
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
	h := hmac.New(sha512.New, keys[currentKeyId].key)
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
	signedToken, err := t.SignedString(keys[currentKeyId].key)
	if err != nil {
		return "", fmt.Errorf("Error in createToken function: %w", err)
	}
	return signedToken, nil
}

func generateNewKey() error {
	newKey := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, newKey)
	if err != nil {
		return fmt.Errorf("Error generating new key: %w", err)
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("Error generating new key ID: %w", err)
	}
	keys[uid.String()] = key{key: newKey, created: time.Now()}
	currentKeyId = uid.String()
	return nil
}

var currentKeyId string
var keys map[string]key

func parseToken(signedToken string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("Invalid signing algorithm")
		}

		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid key ID")
		}

		k, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("Invalid key ID")
		}

		return k.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error in parseToken function: %w", err)
	}
	if !t.Valid {
		return nil, fmt.Errorf("Error in parseToken function: token is not valid.")
	}
	return t.Claims.(*UserClaims), nil
}
