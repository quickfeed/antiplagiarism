package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func pullFiles(labDir string, githubToken string, githubOrg string, studentRepos []string) bool {
	allDownloaded := true

	for _, repo := range studentRepos {
		source := "github.com/" + githubOrg + "/" + repo + ".git"
		fullSource := "https://" + githubToken + ":x-oauth-basic@" + source
		destination := filepath.Join(labDir, githubOrg, repo)

		// Remove any existing files
		err := os.RemoveAll(destination)
		if err != nil {
			fmt.Printf("Error while deleting old contents from %s: %s\n", destination, err)
			allDownloaded = false
			continue
		}

		pullCmd := exec.Command("git")
		pullCmd.Args = []string{"git", "clone", fullSource, destination}
		var sout, serr bytes.Buffer
		pullCmd.Stdout, pullCmd.Stderr = &sout, &serr

		err = pullCmd.Start()
		if err != nil {
			fmt.Printf("Error while pulling from github: %s: %s\n", err, pullCmd.Stderr)
			allDownloaded = false
			continue
		}
		fmt.Printf("Started pulling %s\n", source)

		err = pullCmd.Wait()
		if err != nil {
			fmt.Printf("Error while pulling from github: %s: %s\n", err, pullCmd.Stderr)
			allDownloaded = false
			continue
		}
		fmt.Printf("Finished pulling %s\n", source)
	}

	return allDownloaded
}
