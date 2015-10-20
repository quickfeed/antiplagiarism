package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Where to download the student files
var LabFilesBaseDirectory string

// The location of the Moss script
var MossFqn string

// The location of the JPlag jar file
var JplagFqn string

// Where to store the results
var ResultsDirectory string

// The number of files the code can appear in before it is ignored by Moss
var MossThreshold int

// An integer telling dupl to ignore tokens less than the threshold
var DuplThreshold int

// An integer telling JPlag to ignore tokens less than the threshold
var JplagThreshold int

var defaultMossThreshold = 10
var defaultDuplThreshold = 15
var defaultJplagThreshold = 15

// ReadConfig reads the config.txt file and populates the
// global config variables. The return argument indicates
// whether or not the function was successful.
func ReadConfig() bool {
	file, err := os.Open("config.txt")
	if err != nil {
		fmt.Printf("Could not open the config file. %s\n", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "LAB_FILES_BASE_DIRECTORY") {
			LabFilesBaseDirectory = getLineValueString(line)
		} else if strings.HasPrefix(line, "MOSS_FULLY_QUALIFIED_NAME") {
			MossFqn = getLineValueString(line)
		} else if strings.HasPrefix(line, "JPLAG_FULLY_QUALIFIED_NAME") {
			JplagFqn = getLineValueString(line)
		} else if strings.HasPrefix(line, "RESULTS_DIRECTORY") {
			ResultsDirectory = getLineValueString(line)
		} else if strings.HasPrefix(line, "MOSS_THRESHOLD") {
			MossThreshold = getLineValueInt(line, defaultMossThreshold)
		} else if strings.HasPrefix(line, "DUPL_THRESHOLD") {
			DuplThreshold = getLineValueInt(line, defaultDuplThreshold)
		} else if strings.HasPrefix(line, "JPLAG_THRESHOLD") {
			JplagThreshold = getLineValueInt(line, defaultJplagThreshold)
		} else {
			fmt.Printf("Not expecting %s in the config file.\n", line)
		}
	}

	/*
	fmt.Printf("%s\n", LabFilesBaseDirectory)
	fmt.Printf("%s\n", MossFqn)
	fmt.Printf("%s\n", JplagFqn)
	fmt.Printf("%s\n", ResultsDirectory)
	fmt.Printf("%d\n", MossThreshold)
	fmt.Printf("%d\n", DuplThreshold)
	fmt.Printf("%d\n", JplagThreshold)
	*/

	err = scanner.Err()
	if err != nil {
		fmt.Printf("Error reading the config file. %s\n", err)
		return false
	}

	return true
}

func getLineValueString(line string) string {
	pos := strings.Index(line, "=")
	return line[pos+1:]
}

func getLineValueInt(line string, defaultValue int) int {
	pos := strings.Index(line, "=")
	tempValue, err := strconv.Atoi(line[pos+1:])
	if err != nil {
		fmt.Printf("Error converting config line %s to an integer.\n", line)
		fmt.Printf("Using default value %d.\n", defaultValue)
		return defaultValue
	}
	return tempValue
}
