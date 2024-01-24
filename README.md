# github-actions-metrics

A GitHub app for collecting metrics on GitHub Actions usage. This project is a golang HTTP API for receiving GitHub webhook events, to be used as a GitHub App. The webhook events will be turned into JSON data documents and forwarded to an event/messaging system for further processing.
