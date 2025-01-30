package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
		fmt.Printf("Error: Could not change to repository directory: %v\n", err)
		os.Exit(1)
	}

	if err := exec.Command("git", "fetch", "--prune", "--tags", "--force", "origin").Run(); err != nil {
		fmt.Printf("Error fetching latest tags: %v\n", err)
		os.Exit(1)
	}

	releases, err := findReleases(commitHash)
	if err != nil {
		fmt.Printf("Error finding release tags: %v\n", err)
		os.Exit(1)
	}

	displayReleases(commitHash, releases)
}

func findReleases(commitHash string) ([]Release, error) {
	// Get all version tags sorted by version number
	tags, err := exec.Command("git", "tag", "--sort=version:refname").Output()
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v", err)
	}

	var releases []Release
	for _, tag := range strings.Split(string(tags), "\n") {
		if !strings.HasPrefix(tag, "v") || tag == "" {
			continue
		}

		// Check if our commit is an ancestor of this release
		cmd := exec.Command("git", "merge-base", "--is-ancestor", commitHash, tag)
		isAncestor := cmd.Run() == nil

		releases = append(releases, Release{Tag: tag, IsMatch: isAncestor})
	}

	return releases, nil
}

func displayReleases(commitHash string, releases []Release) {
	fmt.Printf("Results for commit %s:\n\n", commitHash)

	matchFound := false
	for _, release := range releases {
		if release.IsMatch {
			matchFound = true
			fmt.Printf("✓ %s\n", release.Tag)
		}
	}

	if !matchFound {
		fmt.Println("This commit is not in any publicly released version yet.")
	}
}

func compareVersions(v1, v2 string) bool {
	// Remove 'v' prefix
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// First split by + to handle build metadata
	v1Parts := strings.Split(v1, "+")
	v2Parts := strings.Split(v2, "+")

	// Then split the first part by - to separate pre-release
	v1Core := strings.Split(v1Parts[0], "-")
	v2Core := strings.Split(v2Parts[0], "-")

	// Compare major.minor.patch
	parts1 := strings.Split(v1Core[0], ".")
	parts2 := strings.Split(v2Core[0], ".")

	// Compare major.minor.patch numbers
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
			num1, _ := strconv.Atoi(parts1[i])
			num2, _ := strconv.Atoi(parts2[i])
			if num1 != num2 {
					return num1 < num2
			}
	}

	// If core versions are equal, check pre-release versions
	if len(v1Core) != len(v2Core) {
			return len(v1Core) > len(v2Core) // Version with pre-release is older
	}

	// If both have pre-release, compare them
	if len(v1Core) > 1 && len(v2Core) > 1 {
			return v1Core[1] < v2Core[1]
	}

	// If core versions and pre-release are equal, check build metadata
	if len(v1Parts) != len(v2Parts) {
			return len(v1Parts) < len(v2Parts) // Version without metadata is older
	}

	// If both have metadata, compare them
	if len(v1Parts) > 1 && len(v2Parts) > 1 {
			return v1Parts[1] < v2Parts[1]
	}

	return false
}
