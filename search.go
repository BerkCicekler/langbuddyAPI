package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchService struct {
	store Store
}

func NewSearchService(s *Store) *SearchService {
	return &SearchService{store: *s}
}

func (s *SearchService) RegisterRouters(r *mux.Router) {
	r.HandleFunc("/search/", WithJWTAuth(s.handleSearchUser, s.store)).Methods("POST")
}

func (s *SearchService) handleSearchUser(w http.ResponseWriter, r *http.Request) {
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

	users, err := s.store.SearchUser(payload["native"].(string), payload["learning"].(string))

	if err != nil {
		http.Error(w, "MySQL Error", http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, users)
}