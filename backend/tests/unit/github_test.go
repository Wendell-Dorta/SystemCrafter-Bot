package unit_test

import (
	"context"
	"testing"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/github"
)

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		url            string
		expectedOwner  string
		expectedRepo   string
		expectedBranch string
		shouldErr      bool
	}{
		{
			url:            "https://github.com/zen-browser/desktop",
			expectedOwner:  "zen-browser",
			expectedRepo:   "desktop",
			expectedBranch: "",
			shouldErr:      false,
		},
		{
			url:            "github.com/golang/go",
			expectedOwner:  "golang",
			expectedRepo:   "go",
			expectedBranch: "",
			shouldErr:      false,
		},
		{
			url:            "https://github.com/owner/repo/tree/develop",
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedBranch: "develop",
			shouldErr:      false,
		},
		{
			url:            "https://github.com/owner/repo/blob/feature-x/main.go",
			expectedOwner:  "owner",
			expectedRepo:   "repo",
			expectedBranch: "feature-x",
			shouldErr:      false,
		},
		{
			url:       "invalid-url",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		owner, repo, branch, err := github.ParseRepoURL(tt.url)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("Expected error for URL '%s', got nil", tt.url)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for URL '%s': %v", tt.url, err)
		}
		if owner != tt.expectedOwner {
			t.Errorf("Expected owner '%s', got '%s'", tt.expectedOwner, owner)
		}
		if repo != tt.expectedRepo {
			t.Errorf("Expected repo '%s', got '%s'", tt.expectedRepo, repo)
		}
		if branch != tt.expectedBranch {
			t.Errorf("Expected branch '%s', got '%s'", tt.expectedBranch, branch)
		}
	}
}

func TestGitHubClient_FetchPublicRepository_Mock(t *testing.T) {
	client := github.NewGitHubClient()

	// Test invalid URL format
	_, _, _, err := client.FetchPublicRepository(context.Background(), "htt p://bad url")
	if err == nil {
		t.Error("Expected error on malformed URL, got nil")
	}
}
