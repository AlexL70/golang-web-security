package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	//	Base64 encoding is used to encode username and password in "Authorization" HTTP
	//	header when using basic HTTP authentication
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}
