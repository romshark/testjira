package action

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/pkg/errors"
)

func JiraReleased() {

	// JIRA Transition
	transitionName := utils.GetEnv("JIRA_TRANSITION")
	releaseTAG := utils.GetEnv("GITHUB_REF")
	projectID := "ECOMDEV"

	// owner, repo, errExtractRepo := utils.ExtractRepo(utils.GetEnv("GITHUB_REPO"))
	// utils.Must(errors.Wrap(errExtractRepo, "invalid GITHUB_REPO"))

	// init jira client
	jiraClient := client.GetJiraClient()

	project, _, errProject := jiraClient.Project.Get(projectID)
	utils.Must(errors.Wrap(errProject, "unable to load JIRA project"))

	versionID := 0
	for _, version := range project.Versions {

		name := strings.ToLower(version.Name)
		name = strings.ReplaceAll(name, "release", "")
		name = strings.ReplaceAll(name, "hotfix", "")
		name = strings.TrimSpace(name)

		// fmt.Println(version.ID, name, version.ReleaseDate, version.Released)

		if strings.Contains(name, releaseTAG) {
			var errStringConversion error
			versionID, errStringConversion = strconv.Atoi(version.ID)
			utils.Must(errStringConversion)
			break
		}
	}

	if versionID == 0 {
		log.Fatal("JIRA version not found")
	}

	version, _, errVersion := jiraClient.Version.Get(versionID)
	utils.Must(errVersion)

	issues, _, errIssues := jiraClient.Issue.Search("project = ECOMDEV AND fixVersion = "+version.ID, nil)
	utils.Must(errors.Wrap(errIssues, "unable to find JIRA issues with given fix version"))

	success := true
	for _, issue := range issues {

		if issue.Fields.Status.Name != "Merged" {
			success = false
			fmt.Println(issue.Key, "unable to apply transition on issue currently in state:", issue.Fields.Status.Name)
			continue
		}

		// transit jira issue state
		errIssueTransition := utils.JiraIssueTransition(issue.Key, transitionName, jiraClient)
		if errIssueTransition != nil {
			success = true
			fmt.Println(issue.Key, "unable to apply transition on issue currently in state:", issue.Fields.Status.Name, errIssueTransition)
			continue
		}
	}

	if !success {
		fmt.Println("release done, not all tickets have been updated automatically. please mark JIRA release version as released", version.ID, version.Name)
		os.Exit(255)
	}
	fmt.Println("release done, please mark JIRA release version as released", version.ID, version.Name)
}
