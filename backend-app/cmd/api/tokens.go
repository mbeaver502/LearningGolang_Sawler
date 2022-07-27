package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Mock user who is a valid user of this API
var validUser = models.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$DoybxYg.0CDm9.iFCcywUeZkUYB42GakhjGXjDT.H9EEZvX6Uq/.u", // "hunter2"
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// should be a DB query
	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	var claims jwt.Claims

	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "example.com"
	claims.Audiences = []string{"example.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
