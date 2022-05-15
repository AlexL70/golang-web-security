package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	msg := "This is totally fun get hands-on and learning it from the ground up. Thank you for sharing this info with me and helping me."
	encoded := encode(msg)
	fmt.Println(encoded)
	decoded := decode(encoded)
	fmt.Println(decoded)
}

func encode(msg string) string {
	return base64.URLEncoding.EncodeToString([]byte(msg))
}

func decode(encoded string) string {
	result, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		log.Panicln(err)
	}
	return string(result)
}
