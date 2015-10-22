package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"fmt"
	"os"
	"path/filepath"
)

var studentRepos []string
var labNames []string
var labLanguages []int
var getResults bool
var gitHubOrg string
var githubToken string

func main() {
	if !ReadConfig() {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
	}

	// Parse command line args
	if !parseArgs() {
		os.Exit(0)
	}

	if !getResults {
		buildAndRunCommands()
	} else {
		checkAndStoreResults()
	}
}

// buildLabInfo() creates the LabInfo structures from the command
// line arguments to give to the various antiplagiarism packages.
// The return argument is the slice of LabInfo structs.
func buildLabInfo() []common.LabInfo {
	var labInfo []common.LabInfo

	for i := range labNames {
		labInfo = append(labInfo, common.LabInfo{labNames[i], gitHubOrg, labLanguages[i]})
	}

	return labInfo
}

func buildAndRunCommands() bool {
	// TODO: Download files from github using oath token from Autograder

	labInfo := buildLabInfo()

	mossCommands, success := moss.CreateCommands(filepath.Join(LabFilesBaseDirectory, gitHubOrg), MossFqn, labInfo, MossThreshold)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
	} else {
		for _, command := range mossCommands {
			fmt.Printf("%s\n", command)
		}
	}

	duplCommands, success := dupl.CreateCommands(filepath.Join(LabFilesBaseDirectory, gitHubOrg), "", labInfo, DuplThreshold)
	if !success {
		fmt.Printf("Error creating the dupl commands.\n")
	} else {
		for _, command := range duplCommands {
			fmt.Printf("%s\n", command)
		}
	}

	jplagCommands, success := jplag.CreateCommands(filepath.Join(LabFilesBaseDirectory, gitHubOrg), JplagFqn, labInfo, JplagThreshold)
	if !success {
		fmt.Printf("Error creating the JPlag commands.\n")
	} else {
		for _, command := range jplagCommands {
			fmt.Printf("%s\n", command)
		}
	}

	// TODO: Run commands

	return true
}

func checkAndStoreResults() bool {
	// TODO: Check if results are there

	// TODO: Store results

	return true
}
