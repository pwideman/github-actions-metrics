package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v58/github"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"goji.io/pat"
)

func main() {
	// TODO: get config from files (yaml?)
	config := baseapp.HTTPConfig{
		Address: "127.0.0.1",
		Port:    3000,
	}
	loggingConfig := baseapp.LoggingConfig{
		Pretty: true,
		Level:  "trace",
	}

	logger := baseapp.NewLogger(loggingConfig)

	// create a server with default options and no metrics prefix
	server, err := baseapp.NewServer(config, baseapp.DefaultParams(logger, "")...)
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("GH_WEBHOOK_SECRET")
	handle := githubevents.New(secret)
	// add callbacks
	handle.OnWorkflowRunEventAny(
		func(deliveryID string, eventName string, event *github.WorkflowRunEvent) error {
			logger.Info().Msgf("Workflow run event of type %s", *event.WorkflowRun.Status)
			j, _ := json.MarshalIndent(event, "", "  ")
			logger.Trace().Msg(string(j))
			return nil
		},
	)

	// register handlers
	server.Mux().Handle(pat.Post(githubapp.DefaultWebhookRoute), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handle.HandleEventRequest(r)
		if err != nil {
			// TODO: set HTTP response code?
			logger.Error().Err(err).Msg("error handling event request")
		}
	}))

	// start the server (blocking)
	server.Start()
}
