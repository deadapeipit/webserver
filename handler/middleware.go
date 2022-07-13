package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// var secureUser = "david"
// var securePass = "david123"

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if strings.HasPrefix(r.URL.Path, "/users") ||
		// 	strings.HasPrefix(r.URL.Path, "/orders") {
		// 	next.ServeHTTP(w, r)
		// }
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJsonResp(w, StatusError, "bad request")
			return
		}

		// // Hashing the password with the default cost of 10
		// securePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// if err != nil {
		// 	return
		// }
		err := bcrypt.CompareHashAndPassword([]byte(Config.SecurePassword), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJsonResp(w, StatusError, "invalid username / password")
			return
		}

		if username != Config.SecureUser { // || string(securePassword) != Config.SecurePassword {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJsonResp(w, StatusError, "invalid username / password")
			return
		}

		next.ServeHTTP(w, r)
	})
}
