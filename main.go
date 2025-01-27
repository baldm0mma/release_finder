package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Release struct {
	Tag     string
	IsMatch bool
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: release-finder <grafana-repo-path> <commit-hash>")
		os.Exit(1)
	}

	repoPath := os.Args[1]
	commitHash := os.Args[2]

	// Change to the repository directory for git commands
	if err := os.Chdir(repoPath); err != nil {
		fmt.Println("Error: Could not change to repository directory: %v\n", err)
		os.Exit(1)
	}

	if err := exec.Command("git", "fetch", "--prune", "--tags", "--force", "origin").Run(); err != nil {
		fmt.Println("Error fetching latest tags: %v\n", err)
		os.Exit(1)
	}

	releases, err := findReleases(commitHash)
	if err != nil {
		fmt.Println("Error finding release tags: %v\n", err)
		os.Exit(1)
	}

	displayReleases(commitHash, releases)
}

func findReleases(commitHash string) ([]Release, error) {
	// Get all version tags - including those with metadata (e.g. v1.0.0-beta1, v1.0.0-rc1, v1.0.0+security, etc.)
	tagRegex := "v[0-9]*.[0-9]*.[0-9]*[-+]*[a-zA-Z0-9.]*"
	tags, err := exec.Command("git", "tag", "--list", tagRegex).Output()
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v", err)
	}

	var releases []Release
	for _, tag := range strings.Split(string(tags), "\n") {
		if tag == "" {
			continue
		}

		// Check if our commit is an ancestor of this release
		cmd := exec.Command("git", "merge-base", "--is-ancestor", commitHash, tag)
		isAncestor := cmd.Run() == nil

		releases = append(releases, Release{Tag: tag, IsMatch: isAncestor})
	}

	sort.Slice(releases, func(i, j int) bool {
		return compareVersions(releases[i].Tag, releases[j].Tag)
	})

	return releases, nil
}

func displayReleases(commitHash string, releases []Release) {
	fmt.Println("Results for commit %s:\n\n", commitHash)

	matchFound := false
	for _, release := range releases {
		if release.IsMatch {
			matchFound = true
			fmt.Println("✓ %s\n", release.Tag)
		}
	}

	if !matchFound {
		fmt.Println("This commit is not in any publicly released version yet.")
	}
}

func compareVersions(version1, version2 string) {}
