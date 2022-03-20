package handlers

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"net/http"

	database "github.com/cynico/oauth/server/db/mysql"
)

type RegisterHandler struct {
	log hclog.Logger
}

func NewRegisterHandler(l hclog.Logger) *RegisterHandler {
	return &RegisterHandler{
		log: l,
	}
}

// HandleRegisterRequests Dummy function for registering clients hitting the endpoint /register (no credentials are required).
func (rh *RegisterHandler) HandleRegisterRequests(w http.ResponseWriter, r *http.Request) {

	// Preparing statement to insert a new client.
	statement, err := database.Db.Prepare("INSERT INTO Clients (client_id, client_secret, redirect_uri) VALUES (?,?,?)")
	if err != nil {
		rh.log.Error("error preparing statement", "err", err)
		httpInternalServerError(w)
		return
	}

	defer func() {
		err := statement.Close()
		if err != nil {
			rh.log.Error("error closing statement", "err", err)
		}
	}()

	// Randomly generating an id and a secret for the client. (Fixed-length for right now).
	newClient := &Client{
		ClientID:     RandStringRunes(100),
		ClientSecret: RandStringRunes(100),
		RedirectURI:  r.URL.Query().Get("redirect_uri"),
	}

	//_, err = statement.Exec(newClient.ClientID, newClient.ClientSecret, newClient.RedirectURI)
	//sqlErr, ok := err.(*mysql.MySQLError)
	regenerateClientData := func(parameters ...any) {
		rh.log.Debug("changing")
		newClient.ClientID = RandStringRunes(100)
		newClient.ClientSecret = RandStringRunes(100)
		newClient.RedirectURI = "testing"
	}

	err = RetryStatement(statement, regenerateClientData, newClient.ClientID, newClient.ClientSecret, newClient.RedirectURI)

	// Retry inserting if the error you encounter is a duplicate entry error.
	// Maximum number of retries defined in util.go
	//i := 0
	//for (i <= maxInsertRetries) && (err != nil && ok && sqlErr.Number == 1062) {
	//	newClient.ClientID = RandStringRunes(100)
	//	newClient.ClientSecret = RandStringRunes(100)
	//	_, err = statement.Exec(newClient.ClientID, newClient.ClientSecret, newClient.RedirectURI)
	//	i++
	//}

	// If it is not, return an error to the client, and log it.
	if err != nil {
		rh.log.Error("error executing statement", "err", err)
		httpInternalServerError(w)
		return
	}

	// Set the content-type header, and send the data in JSON to the client.
	w.Header().Add("Content-Type", "application/json")
	newEncoder := json.NewEncoder(w)
	err = newEncoder.Encode(newClient)

	if err != nil {
		// Rollback: remove from the db, notify the client, and log the error.
		statement, err := database.Db.Prepare("DELETE FROM Clients WHERE client_id = ?")
		_, _ = statement.Exec(newClient.ClientID)

		httpInternalServerError(w)
		rh.log.Error("error sending the credentials to the client", "err", err)
		return
	}

	rh.log.Debug("A new client registration", "client_id", newClient.ClientID, "client_secret", newClient.ClientSecret)
}
