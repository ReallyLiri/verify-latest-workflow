package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v38/github"
	"github.com/willabides/ezactions"
)

//go:generate go run . -generate

var action = &ezactions.Action{
	Name:        "Verify Latest Workflow",
	Description: "Action to verify the successful status of the latest run of a given workflow",
	Inputs:      []ezactions.ActionInput{},
	Outputs:     []ezactions.ActionOutput{},
	Run:         actionMain,
}

func main() {
	if os.Getenv("MANUAL") != "" {
		_, err := actionMain(nil, nil)
		if err != nil {
			panic(err)
		}
	} else {
		action.Main()
	}
}

var token = os.Getenv("GITHUB_TOKEN")
var owner = os.Getenv("OWNER")
var repository = os.Getenv("REPOSITORY")
var workflow = os.Getenv("WORKFLOW")
var branch = os.Getenv("BRANCH")
var event = os.Getenv("EVENT")

func actionMain(_ map[string]string, _ *ezactions.RunResources) (map[string]string, error) {

	if token == "" || owner == "" || repository == "" || workflow == "" {
		return nil, fmt.Errorf("one or more of the required manual.env vars was not provided")
	}

	log.Printf("working with %v/%v/actions/workflows/%v\n", owner, repository, workflow)

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(oauthClient)

	runs, response, err := githubClient.Actions.ListWorkflowRunsByFileName(ctx, owner, repository, workflow, &github.ListWorkflowRunsOptions{
		Branch: branch,
		Event:  event,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	})

	if runs == nil || err != nil || response.StatusCode != 200 {
		if response != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				waitDuration := time.Now().Sub(response.Rate.Reset.Time)
				return nil, fmt.Errorf("hit rate limit, will need to wait %v sec and retry", waitDuration.Seconds())
			} else {
				return nil, fmt.Errorf("failed to query runs, failed with code %v: %v", response.StatusCode, err)
			}
		} else {
			return nil, fmt.Errorf("failed to query runs: %v", err)
		}
	}

	if len(runs.WorkflowRuns) == 0 {
		return nil, fmt.Errorf("no run results found")
	}

	latestRunStatus := runs.WorkflowRuns[0].GetStatus()
	if latestRunStatus != "completed" {
		return nil, fmt.Errorf("latest run status is %v , not completed", latestRunStatus)
	}

	log.Printf("Latest run verified as completed\n")

	return nil, nil
}
