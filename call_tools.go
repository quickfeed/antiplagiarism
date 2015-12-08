package main

import (
	"bytes"
	"fmt"
	"github.com/autograde/antiplagiarism/common"
	"github.com/autograde/antiplagiarism/dupl"
	"github.com/autograde/antiplagiarism/jplag"
	"github.com/autograde/antiplagiarism/moss"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// buildAndRunCommands builds and runs the commands. The return argument
// indicates whether or not the function was successful. It takes as input
// args, the set of command line arguments, and env, the set of environment variables.
func buildAndRunCommands(args *commandLineArgs, env *envVariables) bool {

	// Pull repositories from github using oath token
	downloadCnt, allDownloaded := pullRepos(env.labDir, args.githubToken, args.githubOrg, args.studentRepos)

	if downloadCnt == 0 {
		fmt.Printf("Failed to download any of the requested repositories. Terminating request.\n")
		return false
	} else if !allDownloaded {
		fmt.Printf("Failed to download all the requested repositories. Proceeding with the rest.\n")
	}

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

	var wg sync.WaitGroup

	// For each tool
	for _, toolCmds := range commands {
		wg.Add(1)

		// Run each command for the tool
		go func(commands []string) {
			defer wg.Done()

			for _, command := range commands {
				//fmt.Printf("%s\n", command)
				executeCommand(&command)
			}
		}(toolCmds)
	}

	// Wait for the tools to finish
	wg.Wait()

	fmt.Printf("Finished all commands. Collecting results.\n")

	// Collect results
	for i := range tools {
		tools[i].SaveResults(args.githubOrg, labInfo, ".", env.resultsDir)
	}

	fmt.Printf("Finished collecting results.\n")

	return true
}

// executeCommand executes a command.
// The return argument indicates whether or not the function was successful.
// It takes as input command, a string containing the command.
func executeCommand(command *string) bool {
	if *command != "" {
		var sout, serr bytes.Buffer
		wd, _ := os.Getwd()

		// Split into command and redirected output
		parts := strings.Split(*command, " > ")

		// Break up the arguments
		cmdArgs := strings.Split(parts[0], " ")

		cmd := exec.Command(cmdArgs[0])
		cmd.Args = cmdArgs
		cmd.Dir = wd
		cmd.Stdout, cmd.Stderr = &sout, &serr

		err := cmd.Start()
		if err != nil {
			fmt.Printf("Error during command %s: %s: %s: %s\n", parts[0], err, sout.Bytes(), serr.Bytes())
			return false
		}
		fmt.Printf("Started %s command.\n", cmdArgs[0])

		err = cmd.Wait()
		if err != nil {
			fmt.Printf("Error during command %s: %s: %s: %s\n", parts[0], err, sout.Bytes(), serr.Bytes())
			return false
		}
		fmt.Printf("Finished %s command.\n", cmdArgs[0])

		// If there was redirected output
		if cmd.Stdout != nil && len(parts) > 1 {
			success := createAndWriteFile(parts[1], &sout)
			return success
		}

		return true
	}
	return false
}

// createAndWriteFile creates a new file and writes the data to it.
// The return argument indicates whether or not the function was successful.
// It takes as input name, the name of the file, and data, the bytes to write.
func createAndWriteFile(name string, data *bytes.Buffer) bool {
	file, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", name, err)
		return false
	}
	defer file.Close()

	file.Write(data.Bytes())
	if err != nil {
		fmt.Printf("Error writing data to file %s: %s\n", name, err)
		return false
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
