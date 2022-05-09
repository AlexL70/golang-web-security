package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)
	http.ListenAndServe(":8080", nil)
}

func encode(w http.ResponseWriter, r *http.Request) {
	p1 := Person{FirstName: "Alex", LastName: "Levinson", Age: 51}
	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println(fmt.Errorf("Error encoding struct: %w", err))
	}
}

func decode(w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(fmt.Errorf("Error decoding request: %w", err))
	}
	fmt.Println(p)
}
