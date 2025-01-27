package main

import ()

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

	
}
