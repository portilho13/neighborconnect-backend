package middleware

import (
	"fmt"
	"net/http"

	"github.com/portilho13/neighborconnect-backend/utils"
)
func Authenticated(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        session, err := utils.Store.Get(r, "session-name")
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