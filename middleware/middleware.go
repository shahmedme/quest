package middleware

import (
	"net/http"

	"../sessions"
)

func AuthRequired(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "cookie-name")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		handlerFunc.ServeHTTP(w, r)
	}
}
