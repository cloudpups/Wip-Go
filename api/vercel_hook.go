package vercel

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v55/github"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	webhookSecretKey := os.Getenv("webhook_secret_key")

	if webhookSecretKey == "" {
		return
	}

	payload, err := github.ValidatePayload(r, []byte(webhookSecretKey))
	if err != nil {
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)

	if err != nil {
		return
	}

	castEvent := event.(*github.PullRequestEvent)

	text := fmt.Sprintf("<h1>The Action was %s</h1>", *castEvent.Action)

	fmt.Fprintf(w, text)
}
