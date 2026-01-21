package middleware

import (
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		next(w, r)
	}
}
