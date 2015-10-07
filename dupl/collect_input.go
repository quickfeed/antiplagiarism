package dupl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"../common"
)

// CreateCommands will create dupl commands to upload the lab files.
// It returns a slice of dupl commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// labsBaseDir, the location of the student directories,
// toolDir, just an empty string for dupl, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling dupl
// to ignore tokens less than the threshold.
func CreateCommands(labsBaseDir string, toolDir string, labs []common.LabInfo, threshold int) ([]string, bool) {
	studentsLabDirs, success := directoryContents(labsBaseDir, labs)
	if !success {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createDuplCommands(studentsLabDirs, labs, threshold)
	if !success {
		fmt.Printf("Error creating the MOSS commands.\n")
		return nil, false
	}

	return commands, true
}

// directoryContents returns a two-dimensional slice with full-path directories
// to send to dupl for evaluation. The first index addresses the specific lab,
// and the second index addresses the specific student. If a student does not
// have the lab directory, the 2D slice will save the directory as an empty
// string. The second return argument indicates whether or not the function was
// successful. directoryContents takes as input baseDir, the location of the
// student directories, and labs, a slice of the labs.
func directoryContents(baseDir string, labs []common.LabInfo) ([][]string, bool) {
	// Try to read the base directory
	contents, err := ioutil.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", baseDir, err)
		return nil, false
	}

	var studentDirs []string
	// Get a list of all the student directories (full path)
	for _, item := range contents {
		if item.IsDir() {
			studentDirs = append(studentDirs, filepath.Join(baseDir, item.Name()))
		}
	}

	studentsLabDirs := make([][]string, len(labs))
	// For each lab
	for i := range studentsLabDirs {
		studentsLabDirs[i] = make([]string, len(studentDirs))

		// For each student
		for j := range studentsLabDirs[i] {
			tempDir := filepath.Join(studentDirs[j], labs[i].Name)
			_, err := ioutil.ReadDir(tempDir)
			if err != nil {
				studentsLabDirs[i][j] = ""
			} else {
				studentsLabDirs[i][j] = tempDir
			}
		}
	}

	return studentsLabDirs, true
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
