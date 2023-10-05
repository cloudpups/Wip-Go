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

	owner := event.GetOrganization()
	repo := event.GetRepo()
	prNum := event.GetPullRequest().GetNumber()
	sha := event.GetPullRequest().GetHead().SHA
	installationID := githubapp.GetInstallationIDFromEvent(&event)

	_, logger := githubapp.PreparePRContext(ctx, installationID, repo, prNum)

	logger.Debug().Msgf("Event action is %s", event.GetAction())

	if !(*event.Action == "labeled" || *event.Action == "unlabeled") {
		logger.Debug().Msgf("Event action is %s, which can be ignored.", event.GetAction())
		return nil
	}

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	name := "WIP"
	conclusion := "success" // "success" or "failure"
	status := "completed"

	labels := event.GetPullRequest().Labels

	labelsToFailOn := []string{"do-not-merge"}

	for _, l := range labels {
		name := l.Name
		markAsFailure := slices.Contains(labelsToFailOn, *name)

		if markAsFailure {
			conclusion = "failure"
			break // No need to continue checking, we found a label that should block the PR!
		}
	}

	if _, _, err := client.Checks.CreateCheckRun(ctx, *owner.Login, *repo.Name, github.CreateCheckRunOptions{
		Conclusion: &conclusion,
		Name:       name,
		Status:     &status,
		HeadSHA:    *sha,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
	}

	return nil
}
