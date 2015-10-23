package moss

import (
	"../common"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO: Cleanup this
var currentDirs1 map[string]bool
var currentDirs2 map[string]bool
var currentLang int

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
		var fileExt []string
		currentLang = labs[i].Language
		// Set language option and file extensions
		if labs[i].Language == common.Golang {
			lOption = "-l java"
			fileExt = append(fileExt, "*.go")
		} else if labs[i].Language == common.Cpp {
			lOption = "-l cc"
			fileExt = append(fileExt, "*.cpp")
			fileExt = append(fileExt, "*.h")
		} else if labs[i].Language == common.C {
			lOption = "-l c"
			fileExt = append(fileExt, "*.c")
			fileExt = append(fileExt, "*.h")
		} else {
			lOption = "-l java"
			fileExt = append(fileExt, "*.java")
		}

		// Start creating the Moss command
		var buf bytes.Buffer
		buf.WriteString(mossFqn + " " + lOption + " " + mOption + " -d")

		// For each student
		for j := range studentsLabDirs[i] {

			// If student has the lab
			if studentsLabDirs[i][j] != "" {
				// TODO: Cleanup this
				currentDirs1 = make(map[string]bool)
				currentDirs2 = make(map[string]bool)
				filepath.Walk(studentsLabDirs[i][j], walkFunction)

				for k := range currentDirs1 {
					buf.WriteString(" " + filepath.Join(k, fileExt[0]))
				}
				for k := range currentDirs2 {
					buf.WriteString(" " + filepath.Join(k, fileExt[1]))
				}

				/*
					// Add all the files with the appropriate extensions
					for k := range fileExt {
						buf.WriteString(" " + filepath.Join(studentsLabDirs[i][j], fileExt[k]))
					}
				*/
			}
		}

		buf.WriteString(" > MOSS." + org + "." + labs[i].Name + ".txt &")

		// Add the Moss command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
}

// TODO: Cleanup this
func walkFunction(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		ind1 := strings.LastIndex(path, ".")
		ind2 := strings.LastIndex(path, "/")
		ext := path[ind1:]
		dir := path[:ind2]

		if currentLang == common.Golang {
			if ext == ".go" {
				currentDirs1[dir] = true
			}
		} else if currentLang == common.Cpp {
			if ext == ".cpp" {
				currentDirs1[dir] = true
			} else if ext == ".h" {
				currentDirs2[dir] = true
			}
		} else if currentLang == common.C {
			if ext == ".c" {
				currentDirs1[dir] = true
			} else if ext == ".h" {
				currentDirs2[dir] = true
			}
		} else {
			if ext == ".java" {
				currentDirs1[dir] = true
			}
		}
	}

	return nil
}
