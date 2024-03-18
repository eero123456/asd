package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ctxUserID = contextKey("ctxuserID")

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == "OPTIONS" {
			return
		}

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ctxUserID, userID))
		next(w, r)

	}

}

func handleNewUser(w http.ResponseWriter) {

	userID, err := addUser()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println("added user", userID)

	token, err := createToken(userID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	SetCookie(w, token)

}

func getUserID(r *http.Request) (uint64, error) {

	existingCookie, err := r.Cookie("token")

	if err != nil {
		return 0, err
	}

	return verifyUser(existingCookie.Value)

}

var jwtKey = []byte("secret-key")

type Claims struct {
	UserID uint64 `json:"userID"`
	jwt.RegisteredClaims
}

func createToken(userID uint64) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyUser(tokenString string) (uint64, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}
	if !tkn.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
