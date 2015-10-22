package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"fmt"
	"os"
)

var studentRepos []string
var labNames []string
var labLanguages []int
var getResults bool
var gitHubOrg string
var githubToken string

func main() {
	// Get the environment variables and parse the command line arguments
	if !GetEnvVar() || !parseArgs() {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
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
		labInfo = append(labInfo, common.LabInfo{labNames[i], labLanguages[i]})
	}

	return labInfo
}

func buildAndRunCommands() bool {
	// TODO: Download files from github using oath token from Autograder

	labInfo := buildLabInfo()
	var tools []common.Tool

	tools = append(tools, moss.Moss{LabsBaseDir: LabFilesBaseDirectory, ToolFqn: MossFqn, Threshold: MossThreshold})
	tools = append(tools, dupl.Dupl{LabsBaseDir: LabFilesBaseDirectory, ToolFqn: "", Threshold: DuplThreshold})
	tools = append(tools, jplag.Jplag{LabsBaseDir: LabFilesBaseDirectory, ToolFqn: JplagFqn, Threshold: JplagThreshold})

	for i := range tools {
		createCommands(tools[i], labInfo)
	}

	// TODO: Run commands

	return true
}

func createCommands(t common.Tool, labs []common.LabInfo) []string {
	commands, success := t.CreateCommands(gitHubOrg, labs)
	if !success {
		fmt.Printf("Error creating the commands.\n")
	} else {
		for _, command := range commands {
			fmt.Printf("%s\n", command)
		}
	}

	return commands
}

func checkAndStoreResults() bool {
	// TODO: Check if results are there

	// TODO: Store results

	return true
}
