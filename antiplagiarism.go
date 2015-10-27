package main

import (
	"fmt"
	"os"
)

type commandLineArgs struct {
	studentRepos []string
	labNames     []string
	labLanguages []int
	getResults   bool
	githubOrg    string
	githubToken  string
}

type envVariables struct {
	// Where to download the student files
	labDir string

	// The location of the Moss script
	mossFqn string

	// The location of the JPlag jar file
	jplagFqn string

	// Where to store the results
	resultsDir string

	// The number of files the code can appear in before it is ignored by Moss
	mossThreshold int

	// An integer telling dupl to ignore tokens less than the threshold
	duplThreshold int

	// An integer telling JPlag to ignore tokens less than the threshold
	jplagThreshold int
}

func main() {
	args := commandLineArgs{}
	env := envVariables{}

	// Get the environment variables and parse the command line arguments
	if !GetEnvVar(&env) || !parseArgs(&args) {
		fmt.Printf("Could not start the anti-plagiarism application.\n")
		os.Exit(0)
	}

	if !args.getResults {
		buildAndRunCommands(&args, &env)
	} else {
		collectResults(&env)
	}
}
