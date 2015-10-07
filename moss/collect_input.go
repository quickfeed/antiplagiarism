package moss

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"

	"../common"
)

// CreateCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// labsBaseDir, the location of the student directories,
// toolDir, the location of the Moss script, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling Moss
// to ignore matches that appear in at least that many files.
func CreateCommands(labsBaseDir string, toolDir string, labs []common.LabInfo, threshold int) ([]string, bool) {
	studentsLabDirs, success := common.DirectoryContents(labsBaseDir, labs)
	if !success {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createMossCommands(toolDir, studentsLabDirs, labs, threshold)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
		return nil, false
	}

	return commands, true
}

// createMossCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. createMossCommands takes as input
// mossDir, the location of the Moss script, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling Moss
// to ignore matches that appear in at least that many files.
func createMossCommands(mossDir string, studentsLabDirs [][]string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	mOption := "-m " + strconv.Itoa(threshold)

	// For each lab
	for i := range studentsLabDirs {
		var lOption string
		var fileExt []string

		// Set language option and file extensions
		if labs[i].Language == common.Golang {
			lOption = "-l java"
			fileExt = append(fileExt, "*.go")
		} else if labs[i].Language == common.Cpp {
			lOption = "-l cc"
			fileExt = append(fileExt, "*.cpp")
			fileExt = append(fileExt, "*.h")
		} else {
			lOption = "-l java"
			fileExt = append(fileExt, "*.java")
		}

		// Start creating the moss command
		var buf bytes.Buffer
		buf.WriteString(filepath.Join(mossDir, "moss") + " " + lOption + " " + mOption + " -d")

		// For each student
		for j := range studentsLabDirs[i] {

			// If student has the lab
			if studentsLabDirs[i][j] != "" {

				// Add all the files with the appropriate extensions
				for k := range fileExt {
					buf.WriteString(" " + filepath.Join(studentsLabDirs[i][j], fileExt[k]))
				}
			}
		}

		buf.WriteString(" > " + labs[i].Name + ".txt &")

		// Add the Moss command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
}
