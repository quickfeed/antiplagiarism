package jplag

import (
	"../common"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// SaveResults saves the data from the specified file. It returns
// whether or not the function was successful. SaveResults takes as input
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
// baseDir, where to save the data, orgName, the name of the organization (class)
// and labName, the name of this lab.
func saveLabResults(fullPathDir string, baseDir string, orgName string, labName string) bool {
	resultsDir := filepath.Join(baseDir, orgName, labName, "jplag")

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
		fmt.Printf("Error reading directory %s: %s\n", fullPathDir, err)
		return false
	}

	// For each file
	for _, info := range fileInfos {
		oldFileNameAndPath := filepath.Join(fullPathDir, info.Name())
		newFileNameAndPath := filepath.Join(resultsDir, info.Name())
		fmt.Printf("Moving file %s\n", info.Name())
		fmt.Printf("	from %s\n", oldFileNameAndPath)
		fmt.Printf("	to %s\n", newFileNameAndPath)

		// Move the file to the results directory.
		os.Rename(oldFileNameAndPath, newFileNameAndPath)
		if err != nil {
			fmt.Printf("Error moving the JPlag output file %s. %s\n", info.Name(), err)
			return false
		}
	}

	// Delete the old folder
	os.RemoveAll(fullPathDir)

	return true
}
