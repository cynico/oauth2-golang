package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"math/rand"
	"net/http"
)

func httpInternalServerError(w http.ResponseWriter) {
	_ = httpJSONResponse(
		&JSONError{
			internalSEMsg,
			internalSEDesc,
			http.StatusInternalServerError},
		w)
}

// httpJSONResponse writes to the given ResponseWrite the json-marshalled response.
func httpJSONResponse(jsonResponse JSONResponse, w http.ResponseWriter) error {
	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(jsonResponse.GetStatusCode())
	return jsonEncoder.Encode(jsonResponse)
}

// RandStringRunes is a function to generate a random string of length n.
func RandStringRunes(n int) string {
	//rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RetryStatement(statement *sql.Stmt, regenerateFunc RegenerateStatementData, parameters ...any) error {

	for _, p := range parameters {
		fmt.Println("parameter: ", p)
	}

	_, err := statement.Exec(parameters...)
	sqlErr, ok := err.(*mysql.MySQLError)

	i := 0
	for (i <= maxInsertRetries) && (err != nil && ok && sqlErr.Number == 1062) {
		regenerateFunc(parameters)
		for _, p := range parameters {
			fmt.Println("parameter: ", p)
		}

		_, err = statement.Exec(parameters...)
		i++
	}

	return err
}
