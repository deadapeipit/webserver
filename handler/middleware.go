package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/jwt") ||
			strings.HasPrefix(r.URL.Path, "/users") {
			next.ServeHTTP(w, r)
			return
		}

		// username, password, ok := r.BasicAuth()
		// if !ok {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	WriteJsonResp(w, StatusError, "bad request")
		// 	return
		// }

		// Check to see if this request can go thru
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			w.WriteHeader(http.StatusForbidden)
			WriteJsonResp(w, StatusError, "FORBIDDEN")
			return
		}

		splitToken := strings.Split(auth, "Bearer ")
		if len(splitToken) != 2 {
			w.WriteHeader(http.StatusForbidden)
			WriteJsonResp(w, StatusError, "FORBIDDEN")
			return
		}

		accessToken := splitToken[1]
		if len(accessToken) == 0 {
			w.WriteHeader(http.StatusForbidden)
			WriteJsonResp(w, StatusError, "FORBIDDEN")
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("signing method invalid")
			}

			return []byte(Config.SecretKey), nil
		})
		if err != nil {
			e, ok := err.(*jwt.ValidationError)
			if !ok || ok && e.Errors&jwt.ValidationErrorIssuedAt == 0 { // Don't report error that token used before issued.
				w.WriteHeader(http.StatusBadRequest)
				WriteJsonResp(w, StatusError, "BAD_REQUEST")
				return
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok { //|| !token.Valid {
			w.WriteHeader(http.StatusBadRequest)
			WriteJsonResp(w, StatusError, "BAD_REQUEST")
			return
		}

		uid := claims["uid"].(string)
		if uid == "" {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJsonResp(w, StatusError, "invalid username / password")
			return
		}

		UserWithRole.Username = uid

		// if uid != Config.SecureUser { // || string(securePassword) != Config.SecurePassword {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	WriteJsonResp(w, StatusError, "invalid username / password")
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
