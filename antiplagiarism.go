package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"fmt"
)

func main() {
	if !ReadConfig() {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
	}

	// TODO: Parse command line args

	// TODO: Download files from github using oath token from Autograder

	// TODO: Remove hard-coded labs
	labs := []common.LabInfo{{"lab1", "DAT500", common.Golang}, {"lab2", "DAT500", common.Java}, {"lab3", "DAT500", common.Cpp}}

	mossCommands, success := moss.CreateCommands(LabFilesBaseDirectory, MossFqn, labs, 4)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
	} else {
		for _, command := range mossCommands {
			fmt.Printf("%s\n", command)
		}
	}

	duplCommands, success := dupl.CreateCommands(LabFilesBaseDirectory, "", labs, 15)
	if !success {
		fmt.Printf("Error creating the dupl commands.\n")
	} else {
		for _, command := range duplCommands {
			fmt.Printf("%s\n", command)
		}
	}

	jplagCommands, success := jplag.CreateCommands(LabFilesBaseDirectory, JplagFqn, labs, 15)
	if !success {
		fmt.Printf("Error creating the JPlag commands.\n")
	} else {
		for _, command := range jplagCommands {
			fmt.Printf("%s\n", command)
		}
	}

	// TODO: Run commands

	// TODO: Store results
}
