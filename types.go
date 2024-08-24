package main

import "time"

type User struct {
	ID                     int64     `json:"id"`
	NotifyToken            string    `json:"notifyToken"`
	UserName               string    `json:"userName"`
	Email                  string    `json:"email"`
	Password               string    `json:"password"`
	NativeLanguage         string    `json:"nativeLanguage"`
	LearningLanguage       string    `json:"learningLanguage"`
	Friends                []string  `json:"friends"`
	ReceivedFriendRequests []string  `json:"receivedFriendRequests"`
	CreatedAt              time.Time `json:"createdAt"`
	AccessToken            string    `json:"accessToken"`
	RefreshToken           string    `json:"refreshToken"`
}
