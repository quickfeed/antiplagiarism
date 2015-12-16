package jplag

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/autograde/antiplagiarism/common"
)

// CreateCommands will create JPlag commands to upload the lab files.
// It returns a slice of JPlag commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// org, the GitHub organization name, and labs, a slice of the labs.
func (j Jplag) CreateCommands(org string, labs []common.LabInfo) ([]string, bool) {
	dir := filepath.Join(j.LabsBaseDir, org)

	// Make sure directory exists
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createJPlagCommands(org, dir, j.ToolFqn, labs, j.Threshold)
	if !success {
		fmt.Printf("Error creating the JPlag commands.\n")
		return nil, false
	}

	return commands, true
}

// createJPlagCommands will create JPlag commands to upload the lab files.
// It returns a slice of JPlag commands.  The second return argument indicates
// whether or not the function was successful. createJPlagCommands takes as input
// org, the GitHub organization name, labsBaseDir, the location of the student directories,
// jplagFqn, the name and path of the JPlag jar, labs, a slice of the labs,
// and threshold, an integer telling JPlag to ignore tokens less than the threshold.
func createJPlagCommands(org string, labsBaseDir string, jplagFqn string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	tOption := " -t " + strconv.Itoa(threshold)

	// For each lab
	for i := range labs {
		var lOption string

		// Set language option
		if labs[i].Language == common.Java {
			lOption = " -l java17"
		} else if labs[i].Language == common.Cpp {
			lOption = " -l c/c++"
		} else if labs[i].Language == common.C {
			lOption = " -l c/c++"
		} else {
			commands = append(commands, "")
			continue
		}

		// Set the results (output) directory
		resultDir := "JPLAG." + org + "." + labs[i].Name
		rOption := " -r " + resultDir

		// Select the subdirectories to compare. In this case, the name of the lab
		sOption := " -S " + labs[i].Name + " -s"

		// Start creating the JPlag command
		var buf bytes.Buffer
		buf.WriteString("java -jar ")
		buf.WriteString(jplagFqn)
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
