package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// pullRepos pulls the requested repositories from github.com.
// It returns a count of the repositories downloaded and true if
// all the repositories were successfully downloaded, otherwise it returns false.
// It takes as input labDir, the directory where github should place the files,
// githubToken, a token authorizing this application to pull from github,
// githubOrg, the organization name containing the student repositories, and
// studentRepos, the names of the student repositories.
func pullRepos(labDir string, githubToken string, githubOrg string, studentRepos []string) (int, bool) {
	downloadCnt := 0
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
		downloadCnt = downloadCnt + 1
	}

	return downloadCnt, allDownloaded
}
