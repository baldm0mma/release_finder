package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type TestRepo struct {
	path string
	t    *testing.T
}

func setupTestRepo(t *testing.T) *TestRepo {
	// Create a temporary directory for the test repository
	tmpDir, err := ioutil.TempDir("", "release-finder-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	repo := &TestRepo{
		path: tmpDir,
		t:    t,
	}

	// Initialize git repository
	repo.git("init")
	repo.git("config", "user.email", "test@example.com")
	repo.git("config", "user.name", "Test User")

	return repo
}

func (r *TestRepo) cleanup() {
	os.RemoveAll(r.path)
}

func (r *TestRepo) git(args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.path
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("Git command failed: git %v: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func (r *TestRepo) createCommit(message string) string {
	// Create a dummy file with unique content
	filename := filepath.Join(r.path, "dummy.txt")
	content := message + "\n"
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		r.t.Fatalf("Failed to write dummy file: %v", err)
	}

	r.git("add", "dummy.txt")
	r.git("commit", "-m", message)

	// Get and return the commit hash
	return r.git("rev-parse", "HEAD")
}

func (r *TestRepo) createTag(tag string) {
	r.git("tag", tag)
}

func TestFindReleases(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.cleanup()

	// Create some commits and tags
	commit1 := repo.createCommit("Initial commit")
	repo.createTag("v1.0.0")

	commit2 := repo.createCommit("Feature commit")
	repo.createTag("v1.1.0")

	commit3 := repo.createCommit("Bugfix commit")
	repo.createTag("v1.1.1")

	commit4 := repo.createCommit("Unreleased commit")

	tests := []struct {
		name            string
		commitHash      string
		expectedTags    []string
		expectInRelease bool
	}{
		{
			name:            "First commit should be in all releases",
			commitHash:      commit1,
			expectedTags:    []string{"v1.0.0", "v1.1.0", "v1.1.1"},
			expectInRelease: true,
		},
		{
			name:            "Middle commit should be in later releases",
			commitHash:      commit2,
			expectedTags:    []string{"v1.1.0", "v1.1.1"},
			expectInRelease: true,
		},
		{
			name:            "Last tagged commit should be in its release",
			commitHash:      commit3,
			expectedTags:    []string{"v1.1.1"},
			expectInRelease: true,
		},
		{
			name:            "Unreleased commit should not be in any release",
			commitHash:      commit4,
			expectedTags:    []string{},
			expectInRelease: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set working directory to repo path
			oldWd, _ := os.Getwd()
			if err := os.Chdir(repo.path); err != nil {
				t.Fatalf("Failed to change working directory: %v", err)
			}
			defer os.Chdir(oldWd)

			releases, err := findReleases(tt.commitHash)
			if err != nil {
				t.Fatalf("findReleases failed: %v", err)
			}

			// Count matching releases
			matchingReleases := []string{}
			for _, release := range releases {
				if release.IsMatch {
					matchingReleases = append(matchingReleases, release.Tag)
				}
			}

			// Verify the number of releases matches expected
			if len(matchingReleases) != len(tt.expectedTags) {
				t.Errorf("Expected %d matching releases, got %d",
					len(tt.expectedTags), len(matchingReleases))
			}

			// Verify each expected tag is present
			for _, expectedTag := range tt.expectedTags {
				found := false
				for _, release := range matchingReleases {
					if release == expectedTag {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find tag %s in matching releases, but didn't",
						expectedTag)
				}
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"v1.0.0", "v1.0.1", true},
		{"v1.0.0", "v2.0.0", true},
		{"v2.0.0", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.0-beta", "v1.0.0", true},
		{"v1.0.0", "v1.0.0-beta", false},
		{"v1.0.0+build", "v1.0.0", false},
		{"v1.0.0", "v1.0.0+build", true},
		{"v1.1.1", "v1.1.1+security", true},
		{"v1.1.1+security", "v1.1.1+security2", true},
		{"v1.1.1+security", "v1.1.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %v; want %v",
					tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}
