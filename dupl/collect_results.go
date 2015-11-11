package dupl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/autograde/antiplagiarism/common"
)

// SaveResults looks for dupl results and saves them. It returns
// whether or not the function was successful. SaveResults takes as input
// org, the GitHub organization name, and labs, a slice of the labs,
// path, where dupl puts the results, and baseDir, where to save the data.
func (m Dupl) SaveResults(org string, labs []common.LabInfo, path string, baseDir string) bool {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", path, err)
		return false
	}

	// For each lab
	for _, lab := range labs {
		fileName := "DUPL." + org + "." + lab.Name + ".html"
		fileNameAndPath := filepath.Join(path, fileName)
		saveLabResults(fileNameAndPath, baseDir, org, lab.Name)
	}

	return true
}

// saveLabResults moves the dupl results to the output directory. It returns
// whether or not the function was successful. SaveResults takes as input
// fileNameAndPath, where dupl saved the results of this lab,
// baseDir, where to save the data, org, the name of the organization (class)
// and labName, the name of this lab.
func saveLabResults(fileNameAndPath string, baseDir string, org string, labName string) bool {
	resultsDir := filepath.Join(baseDir, org, labName, "dupl")
	newFileNameAndPath := filepath.Join(resultsDir, "index.html")

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

	collectPercentages(org, newFileNameAndPath)
	return true
}
