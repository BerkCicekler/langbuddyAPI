package main

import "time"

type User struct {
	ID                     int64     `json:"id"`
	NotifyToken            string    `json:"notifyToken,omitempty"`
	UserName               string    `json:"userName"`
	Email                  string    `json:"email"`
	Password               string    `json:"password,,omitempty"`
	NativeLanguage         string    `json:"nativeLanguage"`
	LearningLanguage       string    `json:"learningLanguage"`
	Friends                []string  `json:"friends"`
	ReceivedFriendRequests []string  `json:"receivedFriendRequests"`
	CreatedAt              time.Time `json:"createdAt,omitempty"`
	AccessToken            string    `json:"accessToken,omitempty"`
	RefreshToken           string    `json:"refreshToken,omitempty"`
}

type Friend struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
}
