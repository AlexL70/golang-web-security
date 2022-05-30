package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/submit", submitEmail)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	isEqual := true
	cEmail := ""
	cCode := ""
	code := ""
	message := "Not logged in"
	xs := strings.SplitN(c.Value, "|", 2)
	if len(xs) == 2 {
		cCode = xs[0]
		cEmail = xs[1]

		code = getCode(cEmail)
		isEqual = hmac.Equal([]byte(code), []byte(cCode))
		if isEqual {
			message = "Logged in"
		}
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>HMAC example</title>
	</head>
	<body>
	<h5> ` + message + `</h5>
	<p> Cookie value: ` + c.Value + `</p>
		<form action="/submit" method="post">
			<input type="email" name="email" value="` + cEmail + `"/>
			<input type="submit" />
		</form>
	</body>
	</html>`
	io.WriteString(w, html)
}

func getCode(msg string) string {
	h := hmac.New(sha256.New, []byte("i love thursdays when it rains 8732 inches"))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func submitEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	code := getCode(email)

	//	hash/digest
	c := http.Cookie{
		Name:  "session",
		Value: code + "|" + email,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
