package action

import (
	"context"
	"log"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/pkg/errors"
)

func PullRequestOpened() {

	owner, repo, err := utils.ExtractRepo(utils.GetEnv("GITHUB_REPO"))
	utils.Must(errors.Wrap(err, "invalid GITHUB_REPO"))

	prKey, errPrKey := utils.GetPR()
	utils.Must(errPrKey)

	// get clients
	ctx := context.Background()
	githubClient := client.GetGithubClientFromEnv(ctx)

	// load PR
	pr, _, err := githubClient.PullRequests.Get(ctx, owner, repo, prKey)
	utils.Must(errors.Wrap(err, "unable to load PR"))

	// extract jira issue key
	issueKey := utils.ExtractJiraIssueKey(pr.GetTitle())

	// init JIRA client
	jiraClient := client.GetJiraClient()

	// prepare update request
	updateVo := struct {
		Fields struct {
			Feature string `json:"customfield_11309"`
			PR      string `json:"customfield_11310"`
		} `json:"fields"`
	}{struct {
		Feature string "json:\"customfield_11309\""
		PR      string "json:\"customfield_11310\""
	}{
		Feature: utils.GetPRHeadBranchURL(pr),
		PR:      pr.GetHTMLURL(),
	}}

	buffer, errBuffer := utils.Encode(updateVo)
	utils.Must(errors.Wrap(errBuffer, "unable to encode updateVo"))

	responseBody, errUpdateIssue := utils.JiraIssueUpdate(issueKey, buffer, jiraClient)
	if errUpdateIssue != nil {
		log.Println("response:", responseBody)
		log.Fatal(errUpdateIssue)
	}

}
