package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	msg := "This is totally fun get hands-on and learning it from the ground up. Thank you for sharing this info with me and helping me."
	password := "My new secure password"
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	logErrorIfAny(err)
	key = key[:16]
	iv := make([]byte, aes.BlockSize) //  initialization vector
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		log.Fatalln(fmt.Errorf("Error randomizing initialization vector: %w", err))
	}
	encrypted, err := encryptMsg(key, []byte(msg), iv)
	logErrorIfAny(err)
	encoded := encode(encrypted)
	fmt.Println(encoded)
	decoded, err := decode(encoded)
	logErrorIfAny(err)
	decrypted, err := encryptMsg(key, decoded, iv)
	logErrorIfAny(err)
	fmt.Println(string(decrypted))
}

func logErrorIfAny(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func encryptMsg(key []byte, msg []byte, iv []byte) ([]byte, error) {
	buff := &bytes.Buffer{}
	sw, err := encryptWriter(buff, key, iv)
	if err != nil {
		return nil, fmt.Errorf("Error getting writer: %w", err)
	}
	defer sw.Close()
	_, err = sw.Write([]byte(msg))
	if err != nil {
		return nil, fmt.Errorf("Error writing message to buffer: %w", err)
	}
	return buff.Bytes(), nil
}

func encryptWriter(wrt io.Writer, key []byte, iv []byte) (*cipher.StreamWriter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error getting new cipher: %w", err)
	}
	s := cipher.NewCTR(block, iv)
	return &cipher.StreamWriter{
		S: s,
		W: wrt,
	}, nil
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
