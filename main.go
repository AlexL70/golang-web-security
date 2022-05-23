package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("sample-file.txt")
	if err != nil {
		log.Fatalln(fmt.Errorf("Error opening file: %w", err))
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatalln(fmt.Errorf("Error copying file to the hash: %w", err))
	}
	fmt.Printf("%x\n%[1]T\n", h.Sum(nil))
}
