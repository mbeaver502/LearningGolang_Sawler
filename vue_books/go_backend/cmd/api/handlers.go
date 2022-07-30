package main

import (
	"books_backend/internal/data"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mozillazg/go-slugify"
)

var staticPath = "./static/"

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

// Login handles logging in a user.
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// authenticate the user
	// look up the user by email (provided by the JSON payload)
	user, err := app.models.User.GetByEmail(creds.Username)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid username or password"))
		return
	}

	// validate the user's provided password
	validPassword, err := user.PasswordMatches(creds.Password)
	if !validPassword || err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid username or password"))
		return
	}

	// make sure user is active
	if user.Active == 0 {
		app.errorJSON(w, errors.New("user is inactive"))
		return
	}

	// if user is valid, generate a token
	token, err := app.models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// save token to database
	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// send back response

	payload := jsonResponse{
		Error:   false,
		Message: "logged in",
		Data: envelope{
			"token": token,
			"user":  user,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}
}

// Logout logs out a user.
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid json"))
		return
	}

	err = app.models.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("failed to delete token"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged out",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// AllUsers gets all users.
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data: envelope{
			"users": all,
		},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// EditUser edits and saves a user.
func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// if ID is 0, then we're adding a new user
	if user.ID == 0 {
		if _, err := user.Insert(); err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}
	} else {
		u, err := app.models.User.GetByID(user.ID)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Active = user.Active

		// update does not change password
		if err = u.Update(); err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}

		// if password != empty string, update password
		if user.Password != "" {
			err := u.ResetPassword(user.Password)
			if err != nil {
				app.errorLog.Println(err)
				app.errorJSON(w, err)
				return
			}
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Changes saved",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// GetUser writes back a requested user.
func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.User.GetByID(userID)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

// DeleteUser deletes a given user.
func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID string `json:"id"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	id, err := strconv.Atoi(requestPayload.ID)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	err = app.models.User.DeleteByID(id)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Changes saved",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// LogUserOutAndSetInactive logs out a user and sets their status to inactive.
func (app *application) LogUserOutAndSetInactive(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.User.GetByID(userID)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	user.Active = 0
	err = user.Update()
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// delete tokens for user
	err = app.models.Token.DeleteTokensForUser(userID)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User logged out and set to inactive",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// ValidateToken will validate a user's token as they navigate to a route.
func (app *application) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	valid := false
	valid, _ = app.models.Token.ValidToken(requestPayload.Token)

	payload := jsonResponse{
		Error: false,
		Data:  valid,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// AllBooks returns all books.
func (app *application) AllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Book.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Data: envelope{
			"books": books,
		},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// OneBook gives back one book.
func (app *application) OneBook(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	book, err := app.models.Book.GetOneBySlug(slug)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Data:  book,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// AllAuthors gives back all authors.
func (app *application) AllAuthors(w http.ResponseWriter, r *http.Request) {
	all, err := app.models.Author.All()
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	type selectData struct {
		Value int    `json:"value"`
		Text  string `json:"text"`
	}

	var results []selectData

	for _, a := range all {
		results = append(results, selectData{
			Value: a.ID,
			Text:  a.AuthorName,
		})
	}

	payload := jsonResponse{
		Error: false,
		Data:  results,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// EditBook edits an existing book or saves a new book.
func (app *application) EditBook(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID              int    `json:"id"`
		Title           string `json:"title"`
		AuthorID        int    `json:"author_id"`
		PublicationYear int    `json:"publication_year"`
		Description     string `json:"description"`
		CoverBase64     string `json:"cover"`
		GenreIDs        []int  `json:"genre_ids"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	book := data.Book{
		ID:              requestPayload.ID,
		Title:           requestPayload.Title,
		AuthorID:        requestPayload.AuthorID,
		PublicationYear: requestPayload.PublicationYear,
		Description:     requestPayload.Description,
		Slug:            slugify.Slugify(requestPayload.Title),
		GenreIDs:        requestPayload.GenreIDs,
	}

	if len(requestPayload.CoverBase64) > 0 {
		// we have a cover -- need to decode and save to FS
		decoded, err := base64.StdEncoding.DecodeString(requestPayload.CoverBase64)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}

		// write image to /static/covers
		// we only accept JPGs, so assume this is JPG
		if err := os.WriteFile(fmt.Sprintf("%s/covers/%s.jpg", staticPath, book.Slug), decoded, 0666); err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}

	}

	if book.ID == 0 {
		// adding a new book
		_, err := app.models.Book.Insert(book)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}
	} else {
		// updating a book
		err := book.Update()
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err)
			return
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Changes saved",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// GetBookByID gets a book by the provided ID.
func (app *application) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	book, err := app.models.Book.GetOneById(id)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Data:  book,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// DeleteBook deletes the requested book.
func (app *application) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	err = app.models.Book.DeleteByID(requestPayload.ID)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Changes saved",
	}

	app.writeJSON(w, http.StatusOK, payload)
}
