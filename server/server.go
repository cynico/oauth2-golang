package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	database "github.com/cynico/oauth/server/db/mysql"
	"github.com/cynico/oauth/server/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

const (
	regexp      = "{[A-Za-z0-9\\-]+}"
	scopeRegexp = "{[0-9]+}"
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "oauth-server",
		Level: hclog.LevelFromString("DEBUG"),
	})

	database.InitDB()
	database.Migrate()

	sm := mux.NewRouter()

	ah := handlers.NewAuthorizingHandler(l)
	rh := handlers.NewRegisterHandler(l)
	mw := handlers.NewMiddleWare(l)
	th := handlers.NewTokenHandler(l)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/o/oauth2/auth", ah.Handle).Queries("client_id", regexp, "state", "{*}", "scope", scopeRegexp)
	getRouter.HandleFunc("/register", rh.HandleRegisterRequests).Queries("redirect_uri", "{*}")

	tokenRouter := sm.Methods(http.MethodPost).Subrouter()
	tokenRouter.Use(mw.BasicAuthMiddleware)
	tokenRouter.Use(mw.ContentTypeMiddleware)
	tokenRouter.HandleFunc("/token", th.GenerateToken)

	s := http.Server{
		Addr:         ":8080",                                                           // configure the bind address
		Handler:      sm,                                                                // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                                   // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                                  // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                                 // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Debug("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server: %s\n", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	l.Debug("Got signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		l.Error("error gracefully shutting down the server", "err", err)
	}
}
