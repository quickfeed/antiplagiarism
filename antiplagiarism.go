package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"flag"
	"fmt"
	"os"
	"strconv"
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
	languages = flag.String(
		"languages",
		"",
		"The list of languages of the labs, delimited by commas. 0 = Java, 1 = Go, 2 = C++",
	)
)

var studentRepos []string
var labNames []string
var labLanguages []int

func usage() {
	flag.PrintDefaults()
	fmt.Printf("Example: ./antiplagiarism -token=0123456789ABCDEF -mainrepo=DAT320 ")
	fmt.Printf("-repos=Student1,Student2,Student3 -labs=LabA,LabB,LabC -languages=0,1,0\n")
}

func main() {
	if !ReadConfig() {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
	}

	// Parse command line args
	if !parseArgs() {
		os.Exit(0)
	}

	// TODO: Download files from github using oath token from Autograder

	labInfo := buildLabInfo()

	mossCommands, success := moss.CreateCommands(LabFilesBaseDirectory, MossFqn, labInfo, MossThreshold)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
	} else {
		for _, command := range mossCommands {
			fmt.Printf("%s\n", command)
		}
	}

	duplCommands, success := dupl.CreateCommands(LabFilesBaseDirectory, "", labInfo, DuplThreshold)
	if !success {
		fmt.Printf("Error creating the dupl commands.\n")
	} else {
		for _, command := range duplCommands {
			fmt.Printf("%s\n", command)
		}
	}

	jplagCommands, success := jplag.CreateCommands(LabFilesBaseDirectory, JplagFqn, labInfo, JplagThreshold)
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

// parseArgs() parses the command line arguments.
// The return argument indicates whether or not the function was successful.
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

	if *languages == "" {
		fmt.Printf("No languages provided.\n")
		return false
	}
	langStr := strings.Split(*languages, ",")

	if len(labNames) != len(langStr) {
		fmt.Printf("The number of labs does not equal the number of languages provided.\n")
		return false
	}

	for _, lang := range langStr {
		langNum, err := strconv.Atoi(lang)
		if err != nil {
			fmt.Printf("Error parsing languages: %v.\n", err)
			return false
		}
		labLanguages = append(labLanguages, langNum)
	}

	return true
}

// buildLabInfo() creates the LabInfo structures from the command
// line arguments to give to the various antiplagiarism packages.
// The return argument is the slice of LabInfo structs.
func buildLabInfo() []common.LabInfo {
	var labInfo []common.LabInfo

	for i := range labNames {
		labInfo = append(labInfo, common.LabInfo{labNames[i], *mainrepo, labLanguages[i]})
	}

	return labInfo
}
