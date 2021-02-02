package action

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_linkJiraIssues(t *testing.T) {
	tests := []struct {
		name string
		diff string
		want string
	}{
		{"standard", "ECOMDEV-123", "[ECOMDEV-123](https://jira.globuswiki.com/browse/ECOMDEV-123)"},
		{"bsi", "BSICRM-321", "[BSICRM-321](https://jira.globuswiki.com/browse/BSICRM-321)"},
		{"extended", "BSICRM-ADD", "BSICRM-ADD"},
		{"text", "Text ECOMDEV-1000 Text", "Text [ECOMDEV-1000](https://jira.globuswiki.com/browse/ECOMDEV-1000) Text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linkJiraIssues(jiraIssueUrlTemplate, tt.diff); got != tt.want {
				t.Errorf("linkJiraIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// This test can be used to generate/re-generate notes
func Test_createReleaseNotes(t *testing.T) {
	token := os.Getenv("GITHUB_KEY")
	tag := os.Getenv("TAG")
	if token == "" || tag == "" {
		t.Skip()
	}

	err := createReleaseNotes("bestbytes", "globus", token, tag)
	assert.NoError(t, err)
}
