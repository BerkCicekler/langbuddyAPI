package main

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func GetUserIdFromRequest(r *http.Request) string {

	claims, err := GetTokenData(GetTokenFromRequest(r))

	if err != nil || claims["userID"] == nil {
		return ""
	}
	return claims["userID"].(string)

}


type ErrorResponse struct {
	Error string `json:"error"`
}

func ContainsSlice(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func RemoveElement(slice []string, valueToRemove string) []string {
	result := []string{}

	for _, value := range slice {
		if value != valueToRemove {
			result = append(result, value)
		}
	}

	return result
}