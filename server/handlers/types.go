package handlers

// JSONResponse is an interface for all json responses: errors and others.
type JSONResponse interface {
	GetStatusCode() int
}

type RegenerateStatementData func (parameters ...any)

// Client is an api client holding its credentials.
type Client struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

// JSONError is a struct for response errors to be marshalled to json.
type JSONError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	StatusCode       int    `json:"-"`
}
func (je * JSONError) GetStatusCode() int {
	return je.StatusCode
}
