// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type PRStatusHandler struct {
	githubapp.ClientCreator
}

func (h *PRStatusHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PRStatusHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse pull request event payload")
	}

	// Requires Palantir --> must seperate
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	//

	// Requires Palantir --> must seperate
	client, err := h.NewInstallationClient(installationID)
	//

	if err != nil {
		return err
	}

	for _, l := range labels {
		name := l.Name
		markAsFailure := slices.Contains(labelsToFailOn, strings.ToLower(*name))

		if markAsFailure {
			conclusion = "failure"
			description = fmt.Sprintf("A blocking label was detected: %s", *name)
			break // No need to continue checking, we found a label that should block the PR!
		}
	}

	if _, _, err := client.Repositories.CreateStatus(ctx, *owner.Login, *repo.Name, *sha, &github.RepoStatus{
		State:       &conclusion,
		Description: &description,
		Context:     &name,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed to create status check")
	}

	var wrapperFunc = func(ctx context.Context, ownerLogin string, repoName string, ref string, status *github.RepoStatus) *error {
		_, _, err := client.Repositories.CreateStatus(ctx, ownerLogin, repoName, ref, &github.RepoStatus{
			State:       &conclusion,
			Description: &description,
			Context:     &name,
		})

		return &err
	}

	CreateStatus(wrapperFunc, ctx, *owner.Login, *repo.Name, *sha, &github.RepoStatus{
		State:       &conclusion,
		Description: &description,
		Context:     &name,
	})

	return nil
}

type StatusCreator func(ctx context.Context, owner string, repo string, ref string, status *github.RepoStatus) *error

func CreateStatus(creator StatusCreator, ctx context.Context, login string, name string, sha string, status *github.RepoStatus) (*github.RepoStatus, error) {
	return creator(ctx, login, name, sha, status)
}

func handlePrLabelEvent(statusCreator StatusCreator, event github.PullRequestEvent) {
	owner := event.GetOrganization()
	repo := event.GetRepo()
	prNum := event.GetPullRequest().GetNumber()
	sha := event.GetPullRequest().GetHead().SHA
	installationId := event.Installation.ID

	if !(*event.Action == "labeled" || *event.Action == "unlabeled") {
		// logger.Debug().Msgf("Event action is %s, which can be ignored.", event.GetAction())
		return
	}

	// TODO: allow name to be configurable so as to account for existing Required Status checks
	name := "WIP"
	conclusion := "success" // "success" or "failure"
	description := "No blocking labels detected!"

	labels := event.GetPullRequest().Labels

	// TODO: make configurable
	// TODO: make sure to create a new array that makes the labelsToFailOn lowercase!!
	labelsToFailOn := []string{"dnm", "do not merge", "do-not-merge", "wip", "work in progress", "work-in-progress"}
}
