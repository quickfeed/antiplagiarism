package dupl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SaveResults looks for dupl results and saves them. It returns
// whether or not the function was successful. SaveResults takes as input
// path, where dupl puts the results, and baseDir, where to save the data.
func (m Dupl) SaveResults(path string, baseDir string) bool {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", path, err)
		return false
	}

	// Regular expression looking for all files starting
	// with DUPL. and ending with .html
	regexStr := "^DUPL.*.html$"
	regex := regexp.MustCompile(regexStr)

	// For each file
	for _, info := range fileInfos {
		fileNameBytes := regex.Find([]byte(info.Name()))

		// If the file name contains the regular expression
		if fileNameBytes != nil {
			fileName := string(fileNameBytes)
			parts := strings.Split(fileName, ".")

			if len(parts) != 4 {
				fmt.Printf("File name %s not in the correct format.\n", fileName)
				continue
			}

			fileNameAndPath := filepath.Join(path, fileName)
			saveLabResults(fileNameAndPath, baseDir, parts[1], parts[2])
		}
	}

	return true
}

// saveLabResults moves the dupl results to the output directory. It returns
// whether or not the function was successful. SaveResults takes as input
// fileNameAndPath, where dupl saved the results of this lab,
// baseDir, where to save the data, orgName, the name of the organization (class)
// and labName, the name of this lab.
func saveLabResults(fileNameAndPath string, baseDir string, orgName string, labName string) bool {
	resultsDir := filepath.Join(baseDir, orgName, labName, "dupl")
	newFileNameAndPath := filepath.Join(resultsDir, "results.html")

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
