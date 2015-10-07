package jplag

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"

	"../common"
)

// CreateCommands will create JPlag commands to upload the lab files.
// It returns a slice of JPlag commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// labsBaseDir, the location of the student directories,
// toolDir, the location of the JPlag jar, labs, a slice of the labs,
// and threshold, an integer telling JPlag to ignore tokens less than the threshold.
func CreateCommands(labsBaseDir string, toolDir string, labs []common.LabInfo, threshold int) ([]string, bool) {
	commands, success := createJPlagCommands(labsBaseDir, toolDir, labs, threshold)
	if !success {
		fmt.Printf("Error creating the JPlag commands.\n")
		return nil, false
	}

	return commands, true
}

// createJPlagCommands will create JPlag commands to upload the lab files.
// It returns a slice of JPlag commands.  The second return argument indicates
// whether or not the function was successful. createJPlagCommands takes as input
// labsBaseDir, the location of the student directories,
// jplagDir, the location of the JPlag jar, labs, a slice of the labs,
// and threshold, an integer telling JPlag to ignore tokens less than the threshold.
func createJPlagCommands(labsBaseDir string, jplagDir string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	tOption := " -t " + strconv.Itoa(threshold)

	// For each lab
	for i := range labs {
		var lOption string

		// Set language option and file extensions
		if labs[i].Language == common.Java {
			lOption = " -l java17"
		} else if labs[i].Language == common.Cpp {
			lOption = " -l c/c++"
		} else {
			continue
		}

		resultDir := filepath.Join(labsBaseDir, "result", labs[i].Class, labs[i].Name)
		rOption := " -r " + resultDir

		sOption := " -s " + labs[i].Name

		// Start creating the JPlag command
		var buf bytes.Buffer
		buf.WriteString("java -jar ")
		buf.WriteString(filepath.Join(jplagDir, "jplag.jar"))
		buf.WriteString(tOption)
		buf.WriteString(lOption)
		buf.WriteString(rOption)
		buf.WriteString(sOption)
		buf.WriteString(" " + labsBaseDir)

		// Add the JPlag command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
}
