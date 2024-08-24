package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserUpdaterService struct {
	store Store
}

func NewUserUpdaterService(s *Store) *UserUpdaterService {
	return &UserUpdaterService{store: *s}
}

func (s *UserUpdaterService) RegisterRouters(r *mux.Router) {
	r.HandleFunc("/user/language", WithJWTAuth(s.handleLanguageSetter, s.store)).Methods("POST")
}

func (s *UserUpdaterService) handleLanguageSetter(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	if payload["native"] == nil || payload["learning"] == nil {
		http.Error(w, "JSON Data is invalid", http.StatusBadRequest)
		return
	}

	claims, err := GetTokenData(GetTokenFromRequest(r))

	if err != nil {
		permissionDenied(w)
		return
	}

	isUpdated, err := s.store.UpdateUserLangs(claims["userID"].(string), payload["native"].(string), payload["learning"].(string))

	if err != nil || !isUpdated {
		log.Print(err)
		http.Error(w, "MySQL Error", http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, "{}")

}
