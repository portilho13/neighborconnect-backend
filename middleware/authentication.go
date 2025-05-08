package middleware

import (
	"fmt"
	"net/http"

	"github.com/portilho13/neighborconnect-backend/utils"
)

func RequireAuthentication(roleRequired string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := utils.Store.Get(r, "user-session")
			if err != nil || session.IsNew {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role, ok := session.Values["role"].(string)
			if !ok || role != roleRequired {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Ex: user_id e email também podem ser verificados aqui se necessário
			_, idOK := session.Values["user_id"]
			_, emailOK := session.Values["email"]
			if !idOK || !emailOK {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}



func AuthenticatedClient(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        session, err := utils.Store.Get(r, "client-session")
        if err != nil {
            fmt.Println("Session error:", err)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if session.IsNew {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        _, userOK := session.Values["user_id"]
        _, emailOK := session.Values["email"]
        
        
        if !userOK || !emailOK {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func AuthenticatedManager(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        session, err := utils.Store.Get(r, "manager-session")
        if err != nil {
            fmt.Println("Session error:", err)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if session.IsNew {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        _, managerOK := session.Values["user_id"]
        _, emailOK := session.Values["email"]
        
        
        if !managerOK || !emailOK {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}