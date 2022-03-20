package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cynico/oauth/client/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	. "github.com/cynico/oauth/client/types"
	browser "github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const (
	state         = "sample-state"
	client_id     = "XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHs"
	client_secret = "bZRjxAwnwekrBEmfdzdcEkXBAkjQZLCtTMtTCoaNatyyiNKAReKJyiXJrscctNswYNsGRussVmaozFZBsbOJiFQGZsnwTKSmVoiG"
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "oauth-client",
		Level: hclog.LevelFromString("DEBUG"),
	})

	sm := mux.NewRouter()
	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Scopes:       []string{"books.read", "books.write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/o/oauth2/auth",
			TokenURL: "http://localhost:8080/o/oauth2/token",
		},
	}

	oauth := OAuth{
		Conf:  conf,
		State: state,
	}

	ch := handlers.NewCallbackHandler(l, &oauth)
	callbackRouter := sm.Methods(http.MethodGet).Subrouter()
	callbackRouter.HandleFunc("/callback", ch.Handle)

	_ = context.Background()

	url := conf.AuthCodeURL("sample-state", oauth2.AccessTypeOffline)
	fmt.Printf("You'll be redirected to the following link: %v\n", url)
	browser.OpenURL(url)

	server := http.Server{
		Addr:         ":7070",
		Handler:      sm,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}),
		ReadTimeout:  5 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		l.Debug("Starting server on port 9090")
		err := server.ListenAndServe()
		if err != nil {
			l.Error("Error starting server: %s\n", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	l.Debug("Got signal: ", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)

}
