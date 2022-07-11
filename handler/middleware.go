package handler

import (
	"net/http"
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
		if username != Config.SecureUser || password != Config.SecurePassword {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJsonResp(w, StatusError, "invalid username / password")
			return
		}
		// fmt.Printf("From config: %s %s \n", Config.SecureUser, Config.SecurePassword)
		// fmt.Printf("From input: %s %s \n", username, password)
		next.ServeHTTP(w, r)
	})
}
