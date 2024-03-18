package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

var secureCookie *securecookie.SecureCookie

var hashKey = []byte("very-secretasdjlkjajdasd0asjda80J=JD=SAJ=DJ=DJKOASJKDO")
var blockKey = []byte("1234abcd----1234")

func initSecureCookie() {
	secureCookie = securecookie.New(hashKey, blockKey)
}

type contextKey string

var latencyCtx = contextKey("myContext")

func setStartTime(r *http.Request) *http.Request {
	timestamp := time.Now().UnixNano()
	return r.WithContext(context.WithValue(r.Context(), latencyCtx, timestamp))
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do some stuff before

		r = setStartTime(r)

		next(w, r)

	}
}

func addLatencyInfo(w http.ResponseWriter, r *http.Request) {

	t1, ok := r.Context().Value(latencyCtx).(int64)
	if !ok {
		return
	}

	duration := time.Now().UnixNano() - t1
	w.Header().Add("X-Duration", fmt.Sprintf("%d", duration))

}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PATCH")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
