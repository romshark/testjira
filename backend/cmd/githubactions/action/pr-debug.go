package action

import (
	"context"
	"fmt"
	"log"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func PullRequestDebug() {

	fmt.Println("GITHUB_REF", utils.GetEnv("GITHUB_REF"))

	prKey, err := utils.ExtractPR(utils.GetEnv("GITHUB_REF"))
	utils.Must(errors.Wrap(err, "invalid GITHUB_REF"))
	fmt.Println("PR", prKey)

	prKey, errPrKey := utils.GetPR()
	utils.Must(errPrKey)
	fmt.Println("PR", prKey)

	// get clients
	ctx := context.Background()
	githubClient := client.GetGithubClientFromEnv(ctx)

	// get PR
	pr, _, errPR := githubClient.PullRequests.Get(ctx, "bestbytes", "globus", prKey)
	utils.Must(errors.Wrap(errPR, "get PR"))
	if pr == nil {
		log.Fatal("PR must not be nil")
	}
	// prURL := pr.GetHTMLURL()

	spew.Dump(pr)
}
