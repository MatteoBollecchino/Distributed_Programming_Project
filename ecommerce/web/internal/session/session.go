package session

// DA USARE I GORILLA SECURE COOKIE

import (
	"net/http"
)

type Session struct {
	UserID string
}

func Set(w http.ResponseWriter, userID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
	})
}
