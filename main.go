package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	msg := "This is totally fun get hands-on and learning it from the ground up. Thank you for sharing this info with me and helping me."
	password := "My new secure password"
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	logErrorIfAny(err)
	key = key[:16]
	encrypted, err := cipherMsg(key, []byte(msg))
	logErrorIfAny(err)
	encoded := encode(encrypted)
	fmt.Println(encoded)
	decoded, err := decode(encoded)
	logErrorIfAny(err)
	decrypted, err := cipherMsg(key, decoded)
	logErrorIfAny(err)
	fmt.Println(string(decrypted))
}

func logErrorIfAny(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func cipherMsg(key []byte, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error getting new cipher")
	}
	s := cipher.NewCTR(block, make([]byte, aes.BlockSize) /*initialization vector*/)
	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}
	defer sw.Close()
	_, err = sw.Write([]byte(msg))
	if err != nil {
		return nil, fmt.Errorf("Error writing message to buffer: %w", err)
	}
	return buff.Bytes(), nil
}

func encode(msg []byte) string {
	return base64.URLEncoding.EncodeToString(msg)
}

func decode(encoded string) ([]byte, error) {
	result, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("Error decoding base64 string: %w", err)
	}
	return result, nil
}
