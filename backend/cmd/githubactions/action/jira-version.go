package action

import (
	"bytes"
	"context"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

func JiraVersion() {

	projectKey := utils.GetEnv("JIRA_PROJECT_KEY")

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

	// get project release versions
	project, _, errProject := jiraClient.Project.Get(projectKey)
	utils.Must(errors.Wrap(errProject, "failed getting jira project"))

	// build release name
	releaseDate := time.Now().Add(time.Hour * 24)
	releaseName := releaseDate.Format("060102") + ".0"

	// iterate release versions and print matching ID
	var version *jira.Version
	for _, v := range project.Versions {
		if strings.Contains(v.Name, releaseName) {
			version = &v
			break
		}
	}

	if version == nil {
		var errVersion error
		version, errVersion = createReleaseVersion(jiraClient, releaseDate, releaseName, project.ID)
		utils.Must(errors.Wrap(errVersion, "unable to create JIRA version"))
	}

	issue, _, errGet := jiraClient.Issue.Get(issueKey, nil)
	utils.Must(errors.Wrap(errGet, "unable to get JIRA ticket"))

	contains := false
	for _, v := range issue.Fields.FixVersions {
		if v.ID == version.ID {
			contains = true
			break
		}
	}

	if !contains {

		// prepare update request
		jsonString := `{"update":{"fixVersions":[{"add":{"name":"` + version.Name + `"}}]}}`

		var buffer io.ReadWriter
		buffer = new(bytes.Buffer)
		_, errWrite := buffer.Write([]byte(jsonString))
		utils.Must(errWrite)

		responseBody, errUpdateIssue := utils.JiraIssueUpdate(issueKey, buffer, jiraClient)
		if errUpdateIssue != nil {
			log.Println("response:", responseBody)
			log.Fatal(errUpdateIssue)
		}

	}

}

func createReleaseVersion(jiraClient *jira.Client, releaseDate time.Time, releaseName, projectID string) (*jira.Version, error) {

	// convert projectID
	id, errConvert := strconv.Atoi(projectID)
	if errConvert != nil {
		return nil, errConvert
	}

	// prepare release version
	version := &jira.Version{
		Name:        releaseName,
		Description: "Release " + releaseName,
		ProjectID:   id,
		ReleaseDate: releaseDate.Format("2006-01-02"),
	}

	// create release version
	version, _, errVersion := jiraClient.Version.Create(version)
	if errVersion != nil {
		return nil, errVersion
	}

	return version, nil
}
