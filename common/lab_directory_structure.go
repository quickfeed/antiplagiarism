package common

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// DirectoryContents returns a two-dimensional slice with full-path directories
// to send to dupl for evaluation. The first index addresses the specific lab,
// and the second index addresses the specific student. If a student does not
// have the lab directory, the 2D slice will save the directory as an empty
// string. The second return argument indicates whether or not the function was
// successful. DirectoryContents takes as input baseDir, the location of the
// student directories, and labs, a slice of the labs.
func DirectoryContents(baseDir string, labs []LabInfo) ([][]string, bool) {
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
