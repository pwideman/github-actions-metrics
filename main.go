package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/bluekeyes/hatpear"
	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v58/github"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"goji.io/pat"

	"github.com/pwideman/github-actions-metrics/internal/config"
	"github.com/pwideman/github-actions-metrics/internal/errors"
)

func main() {
	configFilename := os.Getenv("CONFIG_FILE")
	if configFilename == "" {
		configFilename = "config.yaml"
	}
	config, err := config.ReadConfig(configFilename)
	if err != nil {
		panic(err)
	}

	logger := baseapp.NewLogger(config.Logging)

	// create a server with default options and no metrics prefix
	server, err := baseapp.NewServer(config.Server, baseapp.DefaultParams(logger, "")...)
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("GH_WEBHOOK_SECRET")
	if secret == "" {
		logger.Warn().Msg("GH_WEBHOOK_SECRET not set, will not validate webhook requests")
	}
	handle := githubevents.New(secret)
	// add callbacks
	handle.OnWorkflowRunEventAny(
		func(deliveryID string, eventName string, event *github.WorkflowRunEvent) error {
			workflowRun := event.GetWorkflowRun()
			if workflowRun == nil {
				return errors.NewHTTPError("workflow run event without workflow run", 400)
			}
			logger.Info().Msgf("Workflow run event of type '%s'", *workflowRun.Status)
			j, _ := json.MarshalIndent(event, "", "  ")
			logger.Trace().Msg(string(j))
			return nil
		},
	)

	// register handlers
	server.Mux().Handle(pat.Post(githubapp.DefaultWebhookRoute), hatpear.TryFunc(func(w http.ResponseWriter, r *http.Request) error {
		return handle.HandleEventRequest(r)
	}))

	// start the server (blocking)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
