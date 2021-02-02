package main

import (
	"log"
	"os"

	"testjira/cmd/githubactions/action"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("unexpected args length, please define a single action to execute")
	}

	switch os.Args[1] {
	case "jira-released":
		action.JiraReleased()
	case "jira-transition":
		action.JiraTransition()
	case "jira-version":
		action.JiraVersion()
	case "pr-debug":
		action.PullRequestDebug()
	case "pr-opened":
		action.PullRequestOpened()
	case "pr-reviewed":
		action.PullRequestReviewed()
	case "release-notes":
		action.GenerateReleaseNotesCommand()
	default:
		log.Fatal("unkown github action to execute")
	}

}
