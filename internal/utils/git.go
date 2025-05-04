package utils

import (
	"os/exec"
	"strings"
)

func CheckIfExecInPath(e string) bool {
	_, err := exec.LookPath(e)
	return err == nil
}

func execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), err
}

func IsDirtyGitRepo() bool {
	if !CheckIfExecInPath("git") {
		return false
	}
	isGitRepoStr, _ := execGit("rev-parse", "--is-inside-work-tree")
	if isGitRepoStr != "true" {
		return false
	}
	status, err := execGit("status", "--porcelain")
	if err != nil {
		return false
	}
	return len(status) > 0
}
