package utils

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v29/github"
)

var regex *regexp.Regexp = regexp.MustCompile("^(ECOMDEV|BSICRM)-([0-9]+)")

//ExtractJiraIssueKey from given input
func ExtractJiraIssueKey(input string) string {
	return regex.FindString(input)
}

// GetPR key from env
func GetPR() (int, error) {
	prString := os.Getenv("GITHUB_PR")
	if prString == "" {
		return 0, errors.New("empty GITHUB_PR")
	}
	return strconv.Atoi(prString)
}

// ExtractPR will return the PR number form given GITHUB_REF
func ExtractPR(ref string) (int, error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 4 {
		return 0, errors.New("unexpected ref parts")
	}
	return strconv.Atoi(parts[2])
}

// ExtractRepo will return the repo and owner from GITHUB_REPO
func ExtractRepo(ref string) (string, string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 2 {
		return "", "", errors.New("unable to extract repo and owner from REF")
	}
	return strings.ToLower(parts[0]), strings.ToLower(parts[1]), nil
}

// GetPRHeadBranchURL will return a public HTML URL for PR's head branch
func GetPRHeadBranchURL(pr *github.PullRequest) string {
	baseURL := pr.GetHead().Repo.GetHTMLURL()
	branchName := pr.GetHead().GetRef()
	branchURL := baseURL + "/tree/" + branchName
	return branchURL
}
