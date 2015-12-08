package moss

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/autograde/antiplagiarism/common"
)

// Can't pass these variables to the walk function, so place them here
// Need a separate map for each extension. The maps contain the
// directories which contain files with that extension.
var numberOfExts = 2
var extensionDirs = make([]map[string]bool, numberOfExts)
var fileExts = []string{"", ""}

// CreateCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// org, the GitHub organization name, and labs, a slice of the labs.
func (m Moss) CreateCommands(org string, labs []common.LabInfo) ([]string, bool) {
	dir := filepath.Join(m.LabsBaseDir, org)
	studentsLabDirs, success := common.DirectoryContents(dir, labs)
	if !success {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createMossCommands(org, m.ToolFqn, studentsLabDirs, labs, m.Threshold)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
		return nil, false
	}

	return commands, true
}

// createMossCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. createMossCommands takes as input
// org, the GitHub organization name,
// mossFqn, the name and path of the Moss script, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling Moss
// to ignore matches that appear in at least that many files.
func createMossCommands(org string, mossFqn string, studentsLabDirs [][]string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	mOption := "-m " + strconv.Itoa(threshold)

	// For each lab
	for i := range studentsLabDirs {
		var lOption string

		// Clean the variables
		for j := 0; j < numberOfExts; j++ {
			extensionDirs[j] = make(map[string]bool)
			fileExts[j] = ""
		}

		// Set language option and file extensions
		currentLang := labs[i].Language
		switch currentLang {
		case common.Java:
			lOption = "-l java"
			fileExts[0] = ".java"
		case common.Golang:
			lOption = "-l java"
			fileExts[0] = ".go"
		case common.Cpp:
			lOption = "-l cc"
			fileExts[0] = ".cpp"
			fileExts[1] = ".h"
		case common.C:
			lOption = "-l c"
			fileExts[0] = ".c"
			fileExts[1] = ".h"
		}

		// Start creating the Moss command
		var buf bytes.Buffer
		buf.WriteString(mossFqn + " " + lOption + " " + mOption + " -d")

		// For each student
		for j := range studentsLabDirs[i] {

			// If student has the lab
			if studentsLabDirs[i][j] != "" {
				// Walk through the directories to find which ones
				// contain files with the necessary extensions.
				filepath.Walk(studentsLabDirs[i][j], walkFunction)
			}
		}

		// For all the extensions
		for j := 0; j < numberOfExts; j++ {
			// Add all the directories that contain files with that extension
			for dir := range extensionDirs[j] {
				buf.WriteString(" " + filepath.Join(dir, "*"+fileExts[j]))
			}
		}

		// Finish the command
		buf.WriteString(" > MOSS." + org + "." + labs[i].Name + ".txt")

		// Add the Moss command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
}

// walkFunction checks if the current file's extension matches
// any of the lab's language's extensions. If so, it adds the
// directory to the appropriate map.
func walkFunction(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		extIndex := strings.LastIndex(path, ".")
		if extIndex == -1 {
			return nil
		}
		dirIndex := strings.LastIndex(path, "/")
		ext := path[extIndex:]
		dir := path[:dirIndex]

		// For all the extensions
		for i := 0; i < numberOfExts; i++ {
			// See if that extension matches this file's extension
			if ext == fileExts[i] {
				// It's a match, add the directory to the map
				extensionDirs[i][dir] = true
			}
		}
	}

	return nil
}
