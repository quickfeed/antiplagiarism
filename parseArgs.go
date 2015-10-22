package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	results = flag.Bool(
		"results",
		false,
		"Check and retrieve the results from the various anti-plagiarism applications",
	)
	token = flag.String(
		"token",
		"",
		"The oauth admin token which authorizes this application to connect to github",
	)
	org = flag.String(
		"org",
		"",
		"The name of the GitHub organization",
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

func usage() {
	flag.PrintDefaults()
	fmt.Printf("Example: ./antiplagiarism -token=0123456789ABCDEF -mainrepo=DAT320 ")
	fmt.Printf("-repos=Student1,Student2,Student3 -labs=LabA,LabB,LabC -languages=0,1,0\n")
	fmt.Printf("Example: ./antiplagiarism -results -mainrepo=DAT320 -labs=LabA,LabB,LabC\n")
}

// parseArgs() parses the command line arguments.
// The return argument indicates whether or not the function was successful.
func parseArgs(args *commandLineArgs) bool {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		return false
	}

	if *results {
		args.getResults = true
	} else {
		args.getResults = false
	}

	if *org == "" {
		fmt.Printf("No GitHub organization provided.\n")
		return false
	}
	args.githubOrg = *org

	if *labs == "" {
		fmt.Printf("No lab names provided.\n")
		return false
	}
	args.labNames = strings.Split(*labs, ",")

	if !args.getResults {
		if *token == "" {
			fmt.Printf("No token provided.\n")
			return false
		}
		args.githubToken = *token

		if *repos == "" {
			fmt.Printf("No student repositories provided.\n")
			return false
		}
		args.studentRepos = strings.Split(*repos, ",")

		if *languages == "" {
			fmt.Printf("No languages provided.\n")
			return false
		}
		langStr := strings.Split(*languages, ",")

		if len(args.labNames) != len(langStr) {
			fmt.Printf("The number of labs does not equal the number of languages provided.\n")
			return false
		}

		for _, lang := range langStr {
			langNum, err := strconv.Atoi(lang)
			if err != nil {
				fmt.Printf("Error parsing languages: %v.\n", err)
				return false
			}
			args.labLanguages = append(args.labLanguages, langNum)
		}
	}

	return true
}
