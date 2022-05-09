package main

import (
	"net/http"
)

type person struct {
	First string
}

func main() {
	//	p1 := person{First: "Jenny"}
	//	p2 := person{First: "James"}
	//	xp := []person{p1, p2}
	//	bs, err := json.Marshal(xp)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	fmt.Println("Print JSON:", string(bs))
	//	var xp2 []person
	//	err = json.Unmarshal(bs, &xp2)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	fmt.Println("Unmarshalled struct:", xp2)

	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)
	http.ListenAndServe(":8080", nil)
}

func encode(w http.ResponseWriter, r *http.Request) {

}

func decode(w http.ResponseWriter, r *http.Request) {

}
