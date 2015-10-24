package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	//"bytes"
	"fmt"
	"os"
	//"os/exec"
	//"strings"
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
		checkAndStoreResults()
	}
}

// buildLabInfo() creates the LabInfo structures from the command
// line arguments to give to the various antiplagiarism packages.
// The return argument is the slice of LabInfo structs.
func buildLabInfo(args *commandLineArgs) []common.LabInfo {
	var labInfo []common.LabInfo

	for i := range args.labNames {
		labInfo = append(labInfo, common.LabInfo{args.labNames[i], args.labLanguages[i]})
	}

	return labInfo
}

// buildAndRunCommands builds and runs the commands. The return argument
// indicates whether or not the function was successful.
func buildAndRunCommands(args *commandLineArgs, env *envVariables) bool {

	// Pull repositories from github using oath token
	//if !pullRepos(env.labDir, args.githubToken, args.githubOrg, args.studentRepos) {
	//	fmt.Printf("Failed to download all the requested repositories.\n")
	//}

	labInfo := buildLabInfo(args)
	var tools []common.Tool

	tools = append(tools, moss.Moss{LabsBaseDir: env.labDir, ToolFqn: env.mossFqn, Threshold: env.mossThreshold})
	tools = append(tools, dupl.Dupl{LabsBaseDir: env.labDir, ToolFqn: "", Threshold: env.duplThreshold})
	tools = append(tools, jplag.Jplag{LabsBaseDir: env.labDir, ToolFqn: env.jplagFqn, Threshold: env.jplagThreshold})
	commands := make([][]string, len(tools))

	// Create the commands for each tool
	for i := range tools {
		commands[i] = createCommands(args, tools[i], labInfo)
	}

	// Execute all the commands
	for _, subsetCmds := range commands {
		for _, command := range subsetCmds {
			if command != "" {
				/*
					wd,_ := os.Getwd()
					fmt.Printf("%s\n", wd)
					parts := strings.Split(command, " > ")
					cmdArgs := []string{parts[0]}
					cmd := exec.Command(cmdArgs[0])
					//if cmdArgs[0] == "java" {
					//	cmd.Path = "/usr/bin/java/"
					//} else if cmdArgs[0] == "dupl" {
					//	cmd.Path = "/home/stud/ericfree/go/bin/dupl/"
					//}
					cmd.Args = cmdArgs
					cmd.Dir = wd
					var sout, serr bytes.Buffer
					cmd.Stdout, cmd.Stderr = &sout, &serr
					err := cmd.Start()
					if err != nil {
						fmt.Printf("Error during command %s: %s: %s: %s\n", parts[0], err, cmd.Stdout, cmd.Stderr)
						continue
					}
					fmt.Printf("Started command %s\n", parts[0])

					err = cmd.Wait()
					if err != nil {
						fmt.Printf("Error during command %s: %s: %s: %s\n", parts[0], err, cmd.Stdout, cmd.Stderr)
						continue
					}
					fmt.Printf("Finished command %s\n", parts[0])
				*/
				fmt.Printf("%s\n", command)
			}
		}
	}

	return true
}

// createCommands creates the commands to call/run the anti-plagiarism tool
// for the requested labs. It returns a slice of commands. It takes as input
// t, an object implementing the common.Tool interface, and labs, a slice
// of common.LabInfo objects.
func createCommands(args *commandLineArgs, t common.Tool, labs []common.LabInfo) []string {
	commands, success := t.CreateCommands(args.githubOrg, labs)
	if !success {
		fmt.Printf("Error creating the commands.\n")
		return nil
	}

	return commands
}

func checkAndStoreResults() bool {
	// TODO: Check if results are there

	// TODO: Store results

	return true
}
