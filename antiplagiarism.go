package main

import (
	"fmt"
	"os"
)

// commandLineArgs is a struct containing the command line arguments
type commandLineArgs struct {
	// A slice of student repositories
	studentRepos []string

	// A slice of the lab names to evaluate
	labNames []string

	// A slice of the languages. One for each lab
	labLanguages []int

	// The GitHub organization (class name)
	githubOrg string

	// The admin token for GitHub authorization
	githubToken string

	// The TCP endpoint to listen on.
	endpoint string
}

// envVariables is a struct containing the environment variables
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

	if args.endpoint == "" {
		buildAndRunCommands(&args, &env)
		os.Exit(0)
	}

	startServer(&args, &env)
}
