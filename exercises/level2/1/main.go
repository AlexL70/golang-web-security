package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type MyFormModel struct {
	UserName string
	Passowrd string
}

var usrPwdMap = make(map[string]string)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
		<html lang="en">
			<head>
			    <meta charset="UTF-8">
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
			    <meta name="viewport" content="width=device-width, initial-scale=1.0">
			    <title>Exercise 1 level 2</title>
			</head>
			<body>
				<form action="/register" method="post">
					<label for="userName">User Name:</label>
					<input type="text" placeholder="User Name" id="userName" name="userName" value=""/>
					<br><br>
					<label for="password">Password:</label>
					<input type="password" placeholder="Password" id="password" name="password" value=""/>
					<br><br>
					<input type="submit" />
				</form>
			</body>
		</html>`
	io.WriteString(w, html)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}
	userName := r.FormValue("userName")
	if userName == "" {
		http.Error(w, "empty user name", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}
	encrypted, err := encryptPwd(password)
	if err != nil {
		err = fmt.Errorf("Error enrypting password: %w", err)
		http.Error(w, err.Error(), http.StatusSeeOther)
		return
	}
	usrPwdMap[userName] = encrypted
	for usr, pwd := range usrPwdMap {
		fmt.Printf("User name: %s\nPassword: %s\n", usr, pwd)
	}
	fmt.Println("----------")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func encryptPwd(pwd string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error encrypting password: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(bs)
	return encoded, nil
}
