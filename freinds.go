package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type FriendsService struct {
	store Store
}

func NewFriendsService(s *Store) *FriendsService {
	return &FriendsService{store: *s}
}

func (s *FriendsService) RegisterRouters(r *mux.Router) {
	r.HandleFunc("/friends/accept", WithJWTAuth(s.handleAcceptFriendRequest, s.store)).Methods("POST")
	r.HandleFunc("/friends/reject", WithJWTAuth(s.handleRejectFriendRequest, s.store)).Methods("POST")
	r.HandleFunc("/friends/sendRequest", WithJWTAuth(s.handleSendFriendRequest, s.store)).Methods("POST")
}

func (s *FriendsService) handleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	if payload["targetId"] == nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	targetId := payload["targetId"].(string)
	userId := GetUserIdFromRequest(r)

	targetFriendRequests,err := s.store.GetUsersFriendRequest(targetId)
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}

	if ContainsSlice(targetFriendRequests, userId) {
		WriteJSON(w, http.StatusOK, "{}")
		return
	}

	newTargetFriendRequestsJSON, err := json.Marshal(append(targetFriendRequests, userId))
	if err != nil {
		http.Error(w, "JSON Marshal Error", http.StatusBadRequest)
		return
	}

	userWithNotifyToken,err := s.store.GetNotifyTokenByID(targetId)
	if err!=nil {
		if userWithNotifyToken != nil && userWithNotifyToken.NotifyToken != "" {
			sendNotificationToToken(userWithNotifyToken.NotifyToken, "New Friend Request", "You recieved a new friend request")
		}
	}

	s.store.UpdateUsersFriendRequests(targetId, string(newTargetFriendRequestsJSON))
	WriteJSON(w, http.StatusOK, "{}")
}

func (s *FriendsService) handleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	if payload["targetId"] == nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	targetId := payload["targetId"].(string)
	userId := GetUserIdFromRequest(r)

	if userId == "" {
		permissionDenied(w)
		return
	}

	usersFriendRequests,err := s.store.GetUsersFriendRequest(userId)
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}

	if !ContainsSlice(usersFriendRequests, targetId) {
		return
	}

	newFriendRequestsJSON,err := json.Marshal(RemoveElement(usersFriendRequests, targetId))
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}
	
	usersFriends, err := s.store.GetUsersFriends(userId)
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}

	newUsersFriendsJSON, err := json.Marshal(append(usersFriends, targetId)) 
	if err != nil {
		http.Error(w, "JSON Marshal Error", http.StatusBadRequest)
		return
	}

	targetUserFriends, err := s.store.GetUsersFriends(targetId)
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}
	targetUserFriendsJSON, err :=  json.Marshal(append(targetUserFriends, userId))  
	if err != nil {
		http.Error(w, "JSON Marshal Error", http.StatusBadRequest)
		return
	}

	s.store.UpdateUsersFriendRequests(userId, string(newFriendRequestsJSON))
	s.store.UpdateUsersFriends(userId, string(newUsersFriendsJSON))
	s.store.UpdateUsersFriends(targetId, string(targetUserFriendsJSON))
	WriteJSON(w, http.StatusOK, "{}")
}

func (s *FriendsService) handleRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	if payload["targetId"] == nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	targetId := payload["targetId"].(string)
	userId := GetUserIdFromRequest(r)

	if userId == "" {
		permissionDenied(w)
		return
	}

	usersFriendRequests,err := s.store.GetUsersFriendRequest(userId)
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}

	if !ContainsSlice(usersFriendRequests, targetId) {
		WriteJSON(w, http.StatusOK, "{}")
		return
	}

	newFriendRequestsJSON,err := json.Marshal(RemoveElement(usersFriendRequests, targetId))
	if err != nil {
		http.Error(w, "DB error", http.StatusBadRequest)
		return
	}

	s.store.UpdateUsersFriendRequests(userId, string(newFriendRequestsJSON))
	WriteJSON(w, http.StatusOK, "{}")
}