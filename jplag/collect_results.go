package jplag

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/autograde/antiplagiarism/common"
)

// SaveResults saves the data from the specified file. It returns
// whether or not the function was successful. SaveResults takes as input
// org, the GitHub organization name, and labs, a slice of the labs,
// path, where dupl puts the results, and baseDir, where to save the data.
func (j Jplag) SaveResults(org string, labs []common.LabInfo, path string, baseDir string) bool {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", path, err)
		return false
	}

	// For each lab
	for _, lab := range labs {
		dirName := "JPLAG." + org + "." + lab.Name
		fullPathDir := filepath.Join(path, dirName)
		saveLabResults(fullPathDir, baseDir, org, lab.Name)
	}

	return true
}

// saveLabResults moves the JPlag results to the output directory. It returns
// whether or not the function was successful. SaveResults takes as input
// fullPathDir, where JPlag saved the results of this lab,
// baseDir, where to save the data, org, the name of the organization (class)
// and labName, the name of this lab.
func saveLabResults(fullPathDir string, baseDir string, org string, labName string) bool {
	resultsDir := filepath.Join(baseDir, org, labName, "jplag")

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

	fileInfos, err := ioutil.ReadDir(fullPathDir)
	if err != nil {
		return false
	}

	// For each file
	for _, info := range fileInfos {
		oldFileNameAndPath := filepath.Join(fullPathDir, info.Name())
		newFileNameAndPath := filepath.Join(resultsDir, info.Name())

		// Move the file to the results directory.
		os.Rename(oldFileNameAndPath, newFileNameAndPath)
		if err != nil {
			fmt.Printf("Error moving the JPlag output file %s. %s\n", info.Name(), err)
			return false
		}
	}

	// Delete the old folder
	os.RemoveAll(fullPathDir)

	collectPercentages(org, filepath.Join(resultsDir, "index.html"))

	return true
}
