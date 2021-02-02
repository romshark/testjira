package action

import (
	"context"
	"log"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/pkg/errors"
)

func JiraTransition() {

	// JIRA Transition
	transitionName := utils.GetEnv("JIRA_TRANSITION")

	// PRkey
	prKey, errPrKey := utils.GetPR()
	utils.Must(errPrKey)

	owner, repo, errExtractRepo := utils.ExtractRepo(utils.GetEnv("GITHUB_REPO"))
	utils.Must(errors.Wrap(errExtractRepo, "invalid GITHUB_REPO"))

	// get clients
	ctx := context.Background()
	githubClient := client.GetGithubClientFromEnv(ctx)
	jiraClient := client.GetJiraClient()

	// get PR
	pr, _, errPR := githubClient.PullRequests.Get(ctx, owner, repo, prKey)
	utils.Must(errors.Wrap(errPR, "get PR"))
	if pr == nil {
		log.Fatal("PR must not be nil")
	}

	// transit jira issue state
	issueKey := utils.ExtractJiraIssueKey(pr.GetTitle())
	errIssueTransition := utils.JiraIssueTransition(issueKey, transitionName, jiraClient)
	utils.Must(errors.Wrap(errIssueTransition, "unable to transit JIRA issue"))
}
