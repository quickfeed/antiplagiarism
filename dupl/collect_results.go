package dupl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"../common"
)

// SaveResults saves the data from the specified file. It returns
// whether or not the function was successful. SaveResults takes as input
// fileNameAndPath, the file output by dupl, baseDir, where to save the data,
// and lab, information about the current lab.
func SaveResults(fileNameAndPath string, baseDir string, lab common.LabInfo) bool {
	pos := strings.LastIndex(fileNameAndPath, "/")
	fileName := fileNameAndPath[pos+1:]
	resultsDir := filepath.Join(baseDir, lab.Name)
	newFileNameAndPath := filepath.Join(resultsDir, fileName)

	// Remove old results.
	err := os.RemoveAll(resultsDir)
	if err != nil {
		fmt.Printf("Error removing the old results directory. %s\n", err)
		return false
	}

	// Make the directory.
	err = os.MkdirAll(resultsDir, 0764)
	if err != nil {
		fmt.Printf("Error creating the results directory. %s\n", err)
		return false
	}

	// Move the html file to the results directory.
	os.Rename(fileNameAndPath, newFileNameAndPath)
	if err != nil {
		fmt.Printf("Error moving the dupl html file. %s\n", err)
		return false
	}

	return true
}
