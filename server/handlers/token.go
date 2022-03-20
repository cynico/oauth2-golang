package handlers

import (
	database "github.com/cynico/oauth/server/db/mysql"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strings"
)

type TokenHandler struct {
	log hclog.Logger
}

func NewTokenHandler(l hclog.Logger) *TokenHandler {
	return &TokenHandler{
		log: l,
	}
}

// ValidateReqBody : A function to validate that all the necessary fields exist in a token request body.
func ValidateReqBody(r *http.Request) bool {
	required := []string{"code", "grant_type", "redirect_uri", "client_id"}
	for _, key := range required {
		if _, ok := r.Form[key]; !ok {
			return false
		}
	}
	return true
}

func (t *TokenHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {

	// First, fetching the client ID from the request context (set by the auth middleware).
	clientIDInterface := r.Context().Value(clientIDKey)
	var clientID string

	if c, ok := clientIDInterface.(string); ok {
		clientID = c
	} else {
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}

	// Secondly, fetching the request body.
	err := r.ParseForm()
	if err != nil {
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}

	if !ValidateReqBody(r) {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// Get the code, redirect_uri, grant_type.
	if strings.Join(r.Form["grant_type"], "") != "authorization_code" {

		return
	}

	code := r.Form["code"][0]
	redirectUri := r.Form["redirect_uri"][0]

	// Connect and check with the database.
	statement, err := database.Db.Prepare("SELECT COUNT(*) FROM Codes NATURAL JOIN Clients WHERE Clients.client_id = ? AND code = ? AND redirect_uri = ?")
	if err != nil {
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		t.log.Error("error preparing statement", "err", err)
		return
	}

	// Querying the row.
	var count int
	row := statement.QueryRow(clientID, code, redirectUri)
	if err := row.Scan(&count); err != nil {
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		t.log.Error("error executing statement", "err", err)
		return
	}

	if count != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}
