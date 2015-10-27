package main

import (
	"./common"
	"./dupl"
	"./jplag"
	"./moss"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// collectResults will continually check if results have been returned from the tools
// The return argument indicates whether or not the function was successful. It takes
// as input env, the set of environment variables.
func collectResults(env *envVariables) bool {
	var tools []common.Tool

	tools = append(tools, moss.Moss{LabsBaseDir: env.labDir, ToolFqn: env.mossFqn, Threshold: env.mossThreshold})
	tools = append(tools, dupl.Dupl{LabsBaseDir: env.labDir, ToolFqn: "", Threshold: env.duplThreshold})
	tools = append(tools, jplag.Jplag{LabsBaseDir: env.labDir, ToolFqn: env.jplagFqn, Threshold: env.jplagThreshold})

	var exitSignalReceived = false

	go listenForInterrupt(&exitSignalReceived)

	// While the user has not pressed Ctrl-c
	for exitSignalReceived == false {
		//for i := range tools {
		fmt.Print("Panda\n")
		time.Sleep(time.Second)
		// TODO: Check if results are there

		// TODO: Store results
		//}
	}

	fmt.Printf("Exiting the anitplagiarism application.\n")
	return true
}

// listenForInterrupt listens for the user to tell
// the application to quit, which will allow the
// application to shut down gracefully. It takes as
// input exitSignalReceived, a pointer to a boolean
// flag which is set to true when the user presses Ctrl-c
func listenForInterrupt(exitSignalReceived *bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Wait for interrupt
	<-ch

	fmt.Printf("Recieved command to exit.\n")
	*exitSignalReceived = true
}
