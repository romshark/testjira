package action

import (
	"context"
	"strconv"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/pkg/errors"
)

func PullRequestReviewed() {

	owner, repo, err := utils.ExtractRepo(utils.GetEnv("GITHUB_REPO"))
	utils.Must(errors.Wrap(err, "invalid GITHUB_REPO"))

	prKey, errPrKey := utils.GetPR()
	utils.Must(errPrKey)

	// get clients
	ctx := context.Background()
	githubClient := client.GetGithubClientFromEnv(ctx)

	// load reviews
	reviews, _, err := githubClient.PullRequests.ListReviews(ctx, owner, repo, prKey, nil)
	utils.Must(errors.Wrap(err, "unable to load reviews"))

	isApproved := false
	for _, review := range reviews {
		if *review.State == "APPROVED" && review.GetUser().GetLogin() != "globus-jira" {
			isApproved = true
		}
	}

	if isApproved {
		_, _, err := githubClient.Issues.AddAssignees(ctx, owner, repo, prKey, []string{"globus-jira"})
		utils.Must(errors.Wrap(err, "unable to assign globus-jira user to PR "+strconv.Itoa(prKey)))
	}

}
