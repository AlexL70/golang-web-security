package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)
	http.ListenAndServe(":8080", nil)
}

func encode(w http.ResponseWriter, r *http.Request) {
	p1 := person{First: "Jenny"}
	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		fmt.Println("Encoded bad data:", err)
	}
}

func decode(w http.ResponseWriter, r *http.Request) {
	var p1 person
	err := json.NewDecoder(r.Body).Decode(&p1)
	if err != nil {
		fmt.Println("Decoded bad data:", err)
	}
	log.Println("Decoded data:", p1)
}
