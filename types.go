package main

import "time"

type User struct {
	ID             int64     `json:"id"`
	UserName       string    `json:"userName"`
	Email 		   string 	 `json:"email"` 
	Password       string    `json:"password"`
	NativeLanguage string    `json:"nativeLanguage"`
	LearningLanguage string   `json:"lernningLanguage"`
	Friends 	   []string  `json:"friends"`
	RecievedFriendRequests []string `json:"recievedFriendRequests"`
	CreatedAt      time.Time `json:"createdAt"`
	Token string `json:"token"`
}