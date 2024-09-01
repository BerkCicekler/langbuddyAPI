package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Store interface {
	//Users
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	IsEmailOnUse(email string) (bool, error)
	UpdateUserLangs(id, nativeLang, learningLang string) (bool, error)
	SearchUser(native, learning string) ([]User, error)
	GetNotifyTokenByID(id string) (*User, error)

	// Friends
	GetUsersFriendRequest(userId string) ([]string, error)
	GetUsersFriends(userId string) ([]string, error)
	UpdateUsersFriendRequests(userId, newJsonList string) error
	UpdateUsersFriends(userId, newJsonList string) error
	GetFriendDataFromList(ids []string) ([]Friend, error)
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", u.UserName, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	u.Friends = []string{}
	u.ReceivedFriendRequests = []string{}
	return u, nil
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	var tempFriends, tempRequests string
	err := s.db.QueryRow("SELECT id, username, email, nativeLanguage, learningLanguage, friends, receivedFriendRequests FROM users WHERE id = ?", id).Scan(&u.ID, &u.UserName, &u.Email, &u.NativeLanguage, &u.LearningLanguage, &tempFriends, &tempRequests)
	json.Unmarshal([]byte(tempFriends), &u.Friends)
	json.Unmarshal([]byte(tempRequests), &u.ReceivedFriendRequests)
	return &u, err
}

func (s *Storage) GetNotifyTokenByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT notifyToken FROM users WHERE id = ?", id).Scan(&u.NotifyToken)
	return &u, err
}

func (s *Storage) GetUserByEmail(email string) (*User, error) {
	var u User
	var tempFriends, tempRequests string
	err := s.db.QueryRow("SELECT id, username, email, password, nativeLanguage, learningLanguage, friends, receivedFriendRequests FROM users WHERE email = ?", email).Scan(&u.ID, &u.UserName, &u.Email, &u.Password, &u.NativeLanguage, &u.LearningLanguage, &tempFriends, &tempRequests)
	json.Unmarshal([]byte(tempFriends), &u.Friends)
	json.Unmarshal([]byte(tempRequests), &u.ReceivedFriendRequests)
	return &u, err
}

func (s *Storage) IsEmailOnUse(email string) (bool, error) {
	var isEmailExist bool
	err := s.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM users WHERE email = ?", email).Scan(&isEmailExist)
	return isEmailExist, err
}

func (s *Storage) UpdateUserLangs(userId, nativeLang, learningLang string) (bool, error) {
	rows, err := s.db.Exec("UPDATE users SET nativeLanguage = ?, learningLanguage = ? WHERE id = ?", nativeLang, learningLang, userId)

	if err != nil {
		return false, err
	}
	affected, err := rows.RowsAffected()

	if err != nil || affected == 0 {
		return false, err
	}

	return true, nil
}

func (s *Storage) SearchUser(usersNative, usersLearning string) ([]User, error) {
	var users []User
	var user User

	rows, err := s.db.Query("SELECT id, username FROM users WHERE nativeLanguage = ? AND learningLanguage = ? ORDER BY RAND() LIMIT 10", usersNative, usersLearning)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) GetUsersFriendRequest(userId string) ([]string, error) {
	var jsonString string
	err := s.db.QueryRow("SELECT receivedFriendRequests FROM users WHERE id = ?", userId).Scan(&jsonString)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var friendRequests []string
	err = json.Unmarshal([]byte(jsonString), &friendRequests)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return friendRequests, nil
}

func (s *Storage) GetUsersFriends(userId string) ([]string, error) {
	var jsonString string
	err := s.db.QueryRow("SELECT friends FROM users WHERE id = ?", userId).Scan(&jsonString)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var friends []string
	err = json.Unmarshal([]byte(jsonString), &friends)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return friends, nil
}

func (s *Storage) UpdateUsersFriendRequests(userId string, newJsonList string) error {
	result, err := s.db.Exec("UPDATE users SET receivedFriendRequests = ? WHERE id = ?", newJsonList, userId)
	if err != nil {
		return err
	}

	effected, err := result.RowsAffected()
	if err != nil || effected == 0 {
		return err
	}

	return nil
}

func (s *Storage) UpdateUsersFriends(userId string, newJsonList string) error {
	result, err := s.db.Exec("UPDATE users SET friends = ? WHERE id = ?", newJsonList, userId)
	if err != nil {
		return err
	}

	effected, err := result.RowsAffected()
	if err != nil || effected == 0 {
		return err
	}

	return nil
}

func (s *Storage) GetFriendDataFromList(ids []string) ([]Friend, error) {
	var u []Friend
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	stmt := `SELECT id, username from users where id in (?` + strings.Repeat(",?", len(args)-1) + `)`

	fmt.Println()

	rows, err := s.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user Friend
		if err := rows.Scan(&user.ID, &user.UserName); err != nil {
			return nil, err
		}
		u = append(u, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return u, nil
}
