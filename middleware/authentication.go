package middleware

import (
	//"fmt"
	"net/http"

	"github.com/portilho13/neighborconnect-backend/utils"
)

func RequireAuthentication(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := utils.Store.Get(r, "session")
			if err != nil || session.IsNew {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role, ok := session.Values["role"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authorized := false
			for _, allowed := range allowedRoles {
				if role == allowed {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
