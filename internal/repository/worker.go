package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type GitStore struct {
	repoPath string
	mu       sync.Mutex
}

func NewGitStore(path string) *GitStore {
	return &GitStore{repoPath: path}
}

func (s *GitStore) StartPullWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.SyncRepo()
		case <-ctx.Done():
			return
		}
	}
}

func (s *GitStore) SyncRepo() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println(workerMessage("Checking repo..."))

	if _, err := os.Stat(s.repoPath); os.IsNotExist(err) {
		fmt.Println(workerMessage("Repo not found!"))

		if _, err := os.Stat("data"); os.IsNotExist(err) {
			fmt.Println(workerMessage("Data directory not found!"))
			fmt.Println(workerMessage("Creating data directory..."))
			cmd := exec.Command("mkdir", "data")
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(workerErrorMessage("failed to create directory:", err))
				return
			}
		}

		fmt.Println(workerMessage("Cloning..."))

		cmd := exec.Command("git", "clone", "https://github.com/EXBO-Studio/stalzone-database.git")
		cmd.Dir = "data"

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s\nOutput: %s\n", workerErrorMessage("failed to clone", err), string(output))
			return
		}
	} else {
		fmt.Println(workerMessage("Repo found!"))
	}

	fmt.Println(workerMessage("Checking for updates..."))

	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = s.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(workerErrorMessage("failed to pull", err))
		return
	}

	outStr := string(output)

	if strings.Contains(outStr, "Already up to date.") {
		fmt.Println(workerMessage("No updates found. Skipping database sync."))
		return
	}

	// TODO: Parse directories and sync to PostgresDB
}

func workerMessage(msg string) string {
	return "Git Worker: " + msg
}

func workerErrorMessage(msg string, err error) string {
	return fmt.Sprintf("Git Worker Error: %s:\n%v", msg, err)
}
