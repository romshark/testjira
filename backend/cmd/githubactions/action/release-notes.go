package action

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"testjira/cmd/githubactions/client"
	"testjira/cmd/githubactions/utils"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

const (
	jiraIssueUrlTemplate = "https://jira.globuswiki.com/browse/%s"
)

var (
	jiraIssueRegex = regexp.MustCompile(`(ECOMDEV|BSICRM)-\d+`)
)

func GenerateReleaseNotesCommand() {
	owner, repo, err := utils.ExtractRepo(utils.GetEnv("GITHUB_REPO"))
	utils.Must(err)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN must not be empty")
	}

	tag := os.Getenv("RELEASE_VERSION")
	if tag == "" {
		log.Fatal("RELEASE_VERSION not set")
	}

	err = createReleaseNotes(owner, repo, token, tag)
	utils.Must(err)
}

func createReleaseNotes(owner, repo, token, tag string) error {
	ctx := context.Background()
	githubClient := client.NewGithubClient(ctx, token)
	latestRelease, _, err := githubClient.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return errors.Wrap(err, "could not fetch repo")
	}
	latestTag := latestRelease.GetTagName()

	diff, err := gitDiff(latestTag, tag)
	if err != nil {
		return errors.Wrap(err, "could not get git diff")
	}

	linkedDiff := linkJiraIssues(jiraIssueUrlTemplate, diff)

	release := &github.RepositoryRelease{
		TagName:    github.String(tag),
		Name:       github.String(tag),
		Body:       github.String(linkedDiff),
		Draft:      github.Bool(true),
		Prerelease: github.Bool(false),
	}

	_, _, err = githubClient.Repositories.CreateRelease(ctx, owner, repo, release)
	if err != nil {
		return errors.Wrap(err, "could not create release")
	}

	return nil
}

func gitDiff(from, to string) (string, error) {
	//git log --pretty="%h - %s" '200219.0'...develop
	command := exec.Command("git", "log", `--pretty="%h - %s"`, fmt.Sprintf("%s...%s", from, to))
	out, err := command.CombinedOutput()
	data := string(out)

	if err != nil {
		return data, err
	}

	result := strings.Builder{}
	for _, line := range strings.Split(data, "\n") {
		line := strings.Trim(line, `"`)
		if len(line) == 0 {
			continue
		}
		result.WriteString(line)
		result.WriteString("\n")
	}
	return result.String(), nil
}

func linkJiraIssues(linkTemplate string, diff string) string {
	return jiraIssueRegex.ReplaceAllStringFunc(diff,
		func(s string) string {
			return fmt.Sprintf("[%s](%s)", s, fmt.Sprintf(linkTemplate, s))
		})
}
