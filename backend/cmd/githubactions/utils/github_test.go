package utils

import "testing"

func TestExtractJiraIssueKey(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		input    string
	}{
		{"ECOMDEV", "ECOMDEV-6735", "ECOMDEV-6735 feat: prepare github actions"},
		{"BSICRM", "BSICRM-6735", "BSICRM-6735 feat: prepare github actions"},
		{"empty", "", "FOO-6735 feat: prepare github actions"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issueKey := ExtractJiraIssueKey(tt.input)
			if issueKey != tt.expected {
				t.Fatal("unexpected issue key", issueKey)
			}
		})
	}
}

func TestExtractRepo(t *testing.T) {
	ref := "Codertocat/Hello-World"
	owner, repo, err := ExtractRepo(ref)
	if err != nil {
		t.Fatal(err)
	}
	if owner != "codertocat" && repo != "hello-world" {
		t.Fatal("unexpected repo")
	}
}

func TestExtractPR(t *testing.T) {
	ref := "refs/pull/558/merge"
	pr, err := ExtractPR(ref)
	if err != nil {
		t.Fatal(err)
	}
	if pr != 558 {
		t.Fatal("unexpected pr", pr)
	}
}
