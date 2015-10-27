package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
// indicates whether or not the function was successful. It takes as input
// args, the set of command line arguments, and env, the set of environment variables.
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
			executeCommand(&command)
		}
	}

	//createScript("test2.sh", commands)
	//runScript("./test2.sh")

	return true
}

// createScript is a backup function in case I can't get the commands
// to run from go
func createScript(name string, data [][]string) bool {
	file, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", name, err)
		return false
	}
	defer file.Close()

	file.Chmod(0744)
	file.Write([]byte("#!/bin/sh\n"))

	for _, commands := range data {
		for _, command := range commands {
			if command == "" {
				continue
			}

			file.Write([]byte(command))
			if err != nil {
				fmt.Printf("Error writing data to file %s: %s\n", name, err)
				return false
			}
			file.Write([]byte("\n"))
		}
	}

	return true
}

// runScript is a backup function in case I can't get the commands
// to run from go
func runScript(name string) bool {
	var sout, serr bytes.Buffer
	wd, _ := os.Getwd()
	cmd := exec.Command(name)
	cmd.Args = []string{name}
	cmd.Dir = wd
	cmd.Stdout, cmd.Stderr = &sout, &serr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error during command %s: %s: %s: %s\n", name, err, sout, serr)
		return false
	}
	fmt.Printf("Started %s command.\n", name)

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error during command %s: %s: %s: %s\n", name, err, sout, serr)
		return false
	}
	fmt.Printf("Finished %s command.\n", name)

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
