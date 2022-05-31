package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/submit", submitEmail)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	cEmail := ""
	if err != nil {
		c = &http.Cookie{}
	}

	//	if c.Value != "" {
	//		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	//		parsedToken, _, err := parser.ParseUnverified(c.Value, nil)
	//		if err == nil {
	//			cEmail = parsedToken.Claims.(MyCustomClaims).Email
	//		}
	//	}
	message := "Not logged in"
	//	token, err := getJWT(cEmail)
	//	if err != nil {
	//		message = err.Error()
	//	} else if cEmail != "" && c.Value == token {
	//		message = "Logged in"
	//	}

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

func getJWT(msg string) (string, error) {
	myKey := []byte("i love thursdays when it rains 8732 inches")
	claims := MyCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "alexander.levinson.70@gmail.com",
		},
		Email: msg,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	s, err := token.SignedString(myKey)
	if err != nil {
		return "", fmt.Errorf("Error signing string: %w", err)
	}
	return s, nil
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

	token, err := getJWT(email)
	if err != nil {
		http.Error(w, fmt.Errorf("Error getting token: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	//	hash/digest
	c := http.Cookie{
		Name:  "session",
		Value: token,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
