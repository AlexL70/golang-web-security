package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type MyFormModel struct {
	UserName     string
	Password     string
	PwdEncrypted string
}

var usrPwdMap = make(map[string]string)

const secretKey = "MyVeryLongAndComplicatedSecretKey"

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/success", success)
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
					<h5>Register:</h5>
					<label for="userName">User Name:</label>
					<input type="text" placeholder="User Name" id="userName" name="userName" value=""/>
					<br><br>
					<label for="password">Password:</label>
					<input type="password" placeholder="Password" id="password" name="password" value=""/>
					<br><br>
					<input type="submit" />
				</form>
				<form action="/login" method="post">
					<h5>Login:</h5>
					<label for="userNameL">User Name:</label>
					<input type="text" placeholder="User Name" id="userNameL" name="userName" value=""/>
					<br><br>
					<label for="passwordL">Password:</label>
					<input type="password" placeholder="Password" id="passwordL" name="password" value=""/>
					<br><br>
					<input type="submit" />
				</form>
		</body>
		</html>`
	io.WriteString(w, html)
}

func success(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
		<html lang="en">
			<head>
			    <meta charset="UTF-8">
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
			    <meta name="viewport" content="width=device-width, initial-scale=1.0">
			    <title>Exercise 1 level 2</title>
			</head>
			<body>
				<h2>Login successful!</h2>
			</body>`
	io.WriteString(w, html)
}

func extractCredsFromPostRequest(r *http.Request) (*MyFormModel, error) {
	result := MyFormModel{}
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("Method not allowed: %s. Use POST method.", r.Method)
	}
	result.UserName = r.FormValue("userName")
	if result.UserName == "" {
		return nil, fmt.Errorf("Empty user name.")
	}
	result.Password = r.FormValue("password")
	if result.Password == "" {
		return nil, fmt.Errorf("Empty password.")
	}
	var err error
	result.PwdEncrypted, err = encryptPwd(result.Password)
	if err != nil {
		return nil, fmt.Errorf("Error enrypting password: %w", err)
	}
	return &result, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	creds, err := extractCredsFromPostRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusSeeOther)
		return
	}
	usrPwdMap[creds.UserName] = creds.PwdEncrypted
	for usr, pwd := range usrPwdMap {
		fmt.Printf("User name: %s\nPassword: %s\n", usr, pwd)
	}
	fmt.Println("----------")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	creds, err := extractCredsFromPostRequest(r)
	if err != nil {
		http.Error(w, fmt.Errorf("Login error: %w", err).Error(), http.StatusSeeOther)
		return
	}
	pwd, ok := usrPwdMap[creds.UserName]
	if !ok || !passwordsMatch(creds.Password, pwd) {
		err := fmt.Sprintf("Login failed:\nEntered: %s \nFound: %s", creds.PwdEncrypted, pwd)
		http.Error(w, err, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func passwordsMatch(password, hashB64Encoded string) bool {
	hash, _ := base64.StdEncoding.DecodeString(hashB64Encoded)
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func encryptPwd(pwd string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error encrypting password: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(bs)
	return encoded, nil
}

func getCode(input string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(input))
	if err != nil {
		return nil, fmt.Errorf("Error getting code: %w", err)
	}
	return mac.Sum(nil), nil
}

func createToken(sessionId string) (string, error) {
	code, err := getCode(sessionId)
	if err != nil {
		return "", fmt.Errorf("Error creating token: %w", err)
	}
	signature := base64.StdEncoding.EncodeToString(code)
	return fmt.Sprintf("%s:%s", base64.StdEncoding.EncodeToString([]byte(sessionId)), signature), nil
}

func parseToken(token string) (string, error) {
	const tokenError = "Bad token"
	splitted := strings.Split(token, ":")
	if len(splitted) != 2 {
		return "", errors.New(tokenError)
	}
	sessionId, err := base64.StdEncoding.DecodeString(splitted[0])
	if err != nil {
		return "", fmt.Errorf("%s: %w", tokenError, err)
	}
	sessionIdStr := string(sessionId)
	codeProvided, err := base64.StdEncoding.DecodeString(splitted[1])
	if err != nil {
		return "", fmt.Errorf("%s: %w", tokenError, err)
	}
	code, err := getCode(sessionIdStr)
	if err != nil {
		return "", fmt.Errorf("%s, %w", tokenError, err)
	}
	if !hmac.Equal(codeProvided, code) {
		return "", errors.New(tokenError)
	}
	return sessionIdStr, nil
}
