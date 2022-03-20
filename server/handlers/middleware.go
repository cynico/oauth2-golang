package handlers

import (
	"context"
	"encoding/base64"
	database "github.com/cynico/oauth/server/db/mysql"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strings"
)

type Middleware struct {
	log hclog.Logger
}

func NewMiddleWare(l hclog.Logger) *Middleware {
	return &Middleware{
		log: l,
	}
}

func (m *Middleware) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var count int

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			http.Error(w, "missing valid basic authorization header", http.StatusForbidden)
			return
		}

		dataEncoded := strings.Split(auth, "Basic ")[1]
		dataDecoded := make([]byte, len([]byte(dataEncoded)))

		_, err := base64.StdEncoding.Decode(dataDecoded, []byte(dataEncoded))
		if err != nil {
			http.Error(w, internalErrorMsg, http.StatusInternalServerError)
			m.log.Error("error while decoding id", "error", err)
			return
		}

		client := strings.Split(string(dataDecoded), ":")
		if len(client) != 2 {
			http.Error(w, "invalid authorization header", http.StatusForbidden)
			return
		}
		clientID := client[0]
		clientSecret := client[1]

		// Connect to database.
		statement, err := database.Db.Prepare("SELECT COUNT(*) FROM Clients WHERE client_id = ? AND client_secret = ?")
		if err != nil {
			http.Error(w, "error authorizing", http.StatusInternalServerError)
			m.log.Error("error while preparing secret", "error", err)
			return
		}

		// Executing the statement.
		row := statement.QueryRow(clientID, clientSecret)
		err = row.Scan(&count)
		if err != nil {
			http.Error(w, "error authorizing", http.StatusInternalServerError)
			return
		} else if count != 1 {
			http.Error(w, "Not Authorized", http.StatusForbidden)
			return
		}

		// Attaching clientID to the request for the next handler to access.
		clientIDKey = struct{}{}
		ctx := r.Context()
		ctx = context.WithValue(ctx, clientIDKey, clientID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

func (m *Middleware) ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
