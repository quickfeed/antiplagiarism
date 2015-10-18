package dupl

import (
	"bytes"
	"fmt"
	"strconv"

	"../common"
)

// CreateCommands will create dupl commands to upload the lab files.
// It returns a slice of dupl commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// labsBaseDir, the location of the student directories,
// toolFqn, just an empty string for dupl, labs, a slice of the labs,
// and threshold, an integer telling dupl to ignore tokens less than the threshold.
func CreateCommands(labsBaseDir string, toolFqn string, labs []common.LabInfo, threshold int) ([]string, bool) {
	studentsLabDirs, success := common.DirectoryContents(labsBaseDir, labs)
	if !success {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createDuplCommands(studentsLabDirs, labs, threshold)
	if !success {
		fmt.Printf("Error creating the dupl commands.\n")
		return nil, false
	}

	return commands, true
}

// createDuplCommands will create dupl commands to upload the lab files.
// It returns a slice of dupl commands.  The second return argument indicates
// whether or not the function was successful. createDuplCommands takes as input
// studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling dupl
// to ignore matches that appear in at least that many files.
func createDuplCommands(studentsLabDirs [][]string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	tOption := "-t " + strconv.Itoa(threshold)

	// For each lab
	for i := range studentsLabDirs {

		// Make sure the lab uses go
		if labs[i].Language == common.Golang {
			// Okay
		} else if labs[i].Language == common.Cpp {
			// Can't use dupl on this lab, so empty string
			commands = append(commands, "")
			continue
		} else {
			// Can't use dupl on this lab, so empty string
			commands = append(commands, "")
			continue
		}

		// Start creating the moss command
		var buf bytes.Buffer
		buf.WriteString("dupl" + " " + tOption + " -html")

		// For each student
		for j := range studentsLabDirs[i] {

			// If student has the lab
			if studentsLabDirs[i][j] != "" {
				buf.WriteString(" " + studentsLabDirs[i][j] + "/")
			}
		}

		buf.WriteString(" > " + labs[i].Name + ".html &")

		// Add the dupl command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
}
