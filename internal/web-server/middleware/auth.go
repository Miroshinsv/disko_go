package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	authService "github.com/Miroshinsv/disko_go/internal/auth-service"
)

const (
	AuthHeader           = "X-Token"
	protectedPrefix      = "protected_"
	protectedPrefixAdmin = "protected_admin"
)

func CORSMethodMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")

			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func AuthAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthHeader)
		if token == "" {
			if strings.Contains(mux.CurrentRoute(r).GetName(), protectedPrefixAdmin) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)

			return
		}

		authSrv := authService.GetAuthService()
		dbUser, err := authSrv.GetUserByJWT(token, authService.JWTAuthAudience)

		if dbUser == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !dbUser.IsAdmin() {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", dbUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthHeader)
		if token == "" {
			if strings.Contains(mux.CurrentRoute(r).GetName(), protectedPrefix) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)

			return
		}

		authSrv := authService.GetAuthService()
		dbUser, err := authSrv.GetUserByJWT(token, authService.JWTAuthAudience)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), "user", dbUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
