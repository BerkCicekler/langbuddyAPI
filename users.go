package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersService struct {
	store Store
}

var errEmailRequired = errors.New("email is required")
var errUsernameRequired = errors.New("username is required")
var errPasswordRequired = errors.New("password is required")

func NewUsersService(s *Store) *UsersService {
	return &UsersService{
		store: *s,
	}
}

func (s *UsersService) RegisterRouters(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleCreateUser).Methods("POST")
	r.HandleFunc("/users/login", s.handleLoginUser).Methods("POST")
	r.HandleFunc("/users/refreshToken", s.handleRefreshToken).Methods("POST")
}

func (s *UsersService) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	newTokens, err := CreateNewJWTWithRefreshToken([]byte("SECRET"), payload.RefreshToken)

	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Refresh token invalid"})
		return
	}

	tokenData, err := GetTokenData(newTokens[0])

	if err != nil || tokenData["userID"] == nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Server couldn't find token data"})
		return
	}

	user, err := s.store.GetUserByID(tokenData["userID"].(string))
	if err != nil {
		log.Print(err)
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "User couldn't found"})
		return
	}

	user.AccessToken = newTokens[0]
	user.RefreshToken = newTokens[1]

	WriteJSON(w, http.StatusOK, user)

}

func (s *UsersService) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	emailAlreadyOnUse, err := s.store.IsEmailOnUse(payload.Email)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if emailAlreadyOnUse {
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Email Exist"})
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Password can not be hashed"})
		return
	}
	payload.Password = hashedPassword

	u, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	tokens, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "couldn't create the token"})
		return
	}

	u.AccessToken = tokens[0]
	u.RefreshToken = tokens[1]

	WriteJSON(w, http.StatusCreated, u)

}

func (s *UsersService) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	fmt.Println(payload.Email)
	user, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		log.Print(err)
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Email not exist"})
		return
	}

	err = ControlPassword(payload.Password, user.Password)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid password wrong"})
		return
	}

	tokens, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating token"})
		return
	}

	user.AccessToken = tokens[0]
	user.RefreshToken = tokens[1]

	WriteJSON(w, http.StatusCreated, user)

}

func validateUserPayload(user *User) error {
	if user.UserName == "" {
		return errUsernameRequired
	}

	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(userID int64, w http.ResponseWriter) ([2]string, error) {
	secret := []byte("SECRET")
	tokens, err := CreateJWT(secret, userID)
	if err != nil {
		return [2]string{}, err
	}

	/*
		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: tokens,
		})*/

	return tokens, nil
}
