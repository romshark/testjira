package client

import (
	"testjira/cmd/githubactions/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

const CustomFieldFeatureBranch string = "customfield_11309"
const CustomFieldPR string = "customfield_11310"

// GetJiraClient returns an authenticated jira client
func GetJiraClient() *jira.Client {

	transport := jira.BasicAuthTransport{
		Username: utils.GetEnv("JIRA_USER"),
		Password: utils.GetEnv("JIRA_PASSWORD"),
	}

	client, err := jira.NewClient(transport.Client(), utils.GetEnv("JIRA_BASE_URL"))
	utils.Must(errors.Wrap(err, "jira client init failed"))

	return client
}
