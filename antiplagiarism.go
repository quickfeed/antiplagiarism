package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	token = flag.String(
		"token",
		"",
		"The oauth admin token which authorizes this application to connect to github",
	)
	mainrepo = flag.String(
		"mainrepo",
		"",
		"The main github repository",
	)
	repos = flag.String(
		"repos",
		"",
		"The list of student repositories, delimited by commas",
	)
	labs = flag.String(
		"labs",
		"",
		"The list of labs to check, delimited by commas",
	)
)

var studentRepos []string
var labNames []string

func usage() {
	//fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	//fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Printf("Example: ./antiplagiarism -token=0123456789ABCDEF -mainrepo=DAT320 -repos=Student1,Student2,Student3 -labs=LabA,LabB,LabC\n")
}

func main() {
	if !ReadConfig() {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
	}

	// Parse command line args
	if !parseArgs() {
		os.Exit(0)
	}

	fmt.Printf("Here.\n")

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

// parseArgs parses the command line arguments
func parseArgs() bool {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		return false
	}

	if *token == "" {
		fmt.Printf("No token provided.\n")
		return false
	}

	if *mainrepo == "" {
		fmt.Printf("No main repository provided.\n")
		return false
	}

	if *repos == "" {
		fmt.Printf("No student repositories provided.\n")
		return false
	}
	studentRepos = strings.Split(*repos, ",")

	if *labs == "" {
		fmt.Printf("No lab names provided.\n")
		return false
	}
	labNames = strings.Split(*labs, ",")

	return true
}
