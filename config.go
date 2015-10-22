package main

import (
	"fmt"
	"os"
	"strconv"
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

// GetEnvVar gets the environment variables. The return argument
// indicates whether or not the function was successful.
func GetEnvVar() bool {

	if !getEnvString("LAB_FILES_BASE_DIRECTORY", &LabFilesBaseDirectory) {
		return false
	}
	if !getEnvString("MOSS_FULLY_QUALIFIED_NAME", &MossFqn) {
		return false
	}
	if !getEnvString("JPLAG_FULLY_QUALIFIED_NAME", &JplagFqn) {
		return false
	}
	if !getEnvString("RESULTS_DIRECTORY", &ResultsDirectory) {
		return false
	}
	if !getEnvInt("MOSS_THRESHOLD", &MossThreshold) {
		return false
	}
	if !getEnvInt("DUPL_THRESHOLD", &DuplThreshold) {
		return false
	}
	if !getEnvInt("JPLAG_THRESHOLD", &JplagThreshold) {
		return false
	}

	return true
}

// getEnvString gets an environment variable.
// The return argument indicates whether or not the function was successful.
// It takes as input the name of the environment variable,
// and a pointer where to store the result.
func getEnvString(name string, variable *string) bool {
	*variable = os.Getenv(name)
	if *variable == "" {
		fmt.Printf("%s environment variable not set.\n", name)
		return false
	}
	return true
}

// getEnvInt gets an environment variable and converts it to an integer.
// The return argument indicates whether or not the function was successful.
// It takes as input the name of the environment variable,
// and a pointer where to store the result.
func getEnvInt(name string, variable *int) bool {
	temp := os.Getenv(name)
	if temp == "" {
		fmt.Printf("%s environment variable not set.\n", name)
		return false
	}

	var err error
	*variable, err = strconv.Atoi(temp)
	if err != nil {
		fmt.Printf("Could not convert %s environment variable %v to an integer.\n", name, temp)
		return false
	}

	return true
}
