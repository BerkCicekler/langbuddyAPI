package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateNewJWTWithRefreshToken(secret []byte, refreshToken string) ([2]string, error) {
	token, err := validateJWT(refreshToken)
	if err != nil {
		log.Printf("failed to validate token: %v", err)
		return [2]string{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["expiresAt"].(float64); ok {
			now := time.Now().Unix()

			if now > int64(exp) {
				return [2]string{}, fmt.Errorf("token has expired")
			}

			userIDStr, ok := claims["userID"].(string)
			if !ok {
				return [2]string{}, fmt.Errorf("userID is not a string")
			}

			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				return [2]string{}, fmt.Errorf("cannot convert userID to int64: %v", err)
			}

			// success
			return CreateJWT(secret, userID)

		} else {
			return [2]string{}, fmt.Errorf("token does not have an expiresAt field")
		}
	} else {
		return [2]string{}, fmt.Errorf("invalid token")
	}
}

func CreateJWT(secret []byte, userID int64) ([2]string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 2).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return [2]string{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return [2]string{}, err
	}

	return [2]string{accessTokenString, refreshTokenString}, err
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ControlPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("SECRET"), nil
	})
}

func GetTokenData(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}
