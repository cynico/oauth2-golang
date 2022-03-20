package handlers

import (
	"fmt"
	"net/http"

	database "github.com/cynico/oauth/server/db/mysql"
	"github.com/hashicorp/go-hclog"
)

var RegisteredClients map[string]map[string]string

const internalErrorMsg = "Internal server error"

func init() {
	RegisteredClients = make(map[string]map[string]string)
	RegisteredClients["sample-client-id"] = make(map[string]string)
	RegisteredClients["sample-client-id"]["redirect_uri"] = "http://localhost:7070/callback"
}

type AuthorizingHandler struct {
	log hclog.Logger
}

func NewAuthorizingHandler(l hclog.Logger) *AuthorizingHandler {
	return &AuthorizingHandler{
		log: l,
	}
}

func (a *AuthorizingHandler) Handle(w http.ResponseWriter, r *http.Request) {

	var count int
	var redirectURI string

	// First getting all the query parameters.
	clientID := r.URL.Query().Get("client_id")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")
	scope := r.URL.Query().Get("scope")

	// First checking if a client with this id is registered.
	statement, err := database.Db.Prepare("SELECT COUNT(*), redirect_uri FROM Clients WHERE client_id = ? GROUP BY (client_id)")
	if err != nil {
		a.log.Error("error preparing statement", "error", err)
		httpInternalServerError(w)
		return
	}

	row := statement.QueryRow(clientID)
	err = row.Scan(&count, &redirectURI)
	if err != nil {
		a.log.Error("error scanning row", "error", err)
		httpInternalServerError(w)
		return
	}

	if count != 1 {
		_ = httpJSONResponse(
			&JSONError{
				"unrecognized_client",
				"Unrecognized client. Send a valid client_id.",
				http.StatusForbidden,
			},
			w)
		return
	}

	a.log.Debug(
		"received an authorizing request",
		"client_id",
		clientID,
		"response_type",
		responseType,
		"state",
		state,
		"scope",
		scope,
	)

	// Generating a random code, and inserting it into the database.
	code := RandStringRunes(100)
	statement, err = database.Db.Prepare("INSERT INTO Codes (client_id, code) VALUES (?,?)")
	if err != nil {
		a.log.Error("error preparing statement", "error", err)
		httpInternalServerError(w)
		return
	}

	_, err = statement.Exec(clientID, code)
	if err != nil {
		a.log.Error("error inserting code", "error", err)
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}

	// Constructing the URL to redirect the user to.
	URL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
	http.Redirect(w, r, URL, http.StatusSeeOther)

}
