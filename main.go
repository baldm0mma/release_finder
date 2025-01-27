package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
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

	// Split version from metadata
	v1Parts := strings.FieldsFunc(v1, func(r rune) bool {
		return r == '-' || r == '+'
	})
	v2Parts := strings.FieldsFunc(v2, func(r rune) bool {
		return r == '-' || r == '+'
	})

	// Compare core versions first
	parts1 := strings.Split(v1Parts[0], ".")
	parts2 := strings.Split(v2Parts[0], ".")

	// Compare major.minor.patch
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])
		if num1 != num2 {
			return num1 < num2
		}
	}

	// If core versions are equal, version with metadata is considered newer
	// unless it's a pre-release (with '-')
	if len(v1Parts) != len(v2Parts) {
		v1HasPreRelease := strings.Contains(v1, "-")
		v2HasPreRelease := strings.Contains(v2, "-")

		if v1HasPreRelease != v2HasPreRelease {
			return v1HasPreRelease // Pre-release version is older
		}
	}

	return len(parts1) < len(parts2)
}
