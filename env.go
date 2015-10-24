package main

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnvVar gets the environment variables. The return argument
// indicates whether or not the function was successful.
func GetEnvVar(env *envVariables) bool {

	if !getEnvString("LAB_FILES_BASE_DIRECTORY", &env.labDir) {
		return false
	}
	if !getEnvString("MOSS_FULLY_QUALIFIED_NAME", &env.mossFqn) {
		return false
	}
	if !getEnvString("JPLAG_FULLY_QUALIFIED_NAME", &env.jplagFqn) {
		return false
	}
	if !getEnvString("RESULTS_DIRECTORY", &env.resultsDir) {
		return false
	}
	if !getEnvInt("MOSS_THRESHOLD", &env.mossThreshold) {
		return false
	}
	if !getEnvInt("DUPL_THRESHOLD", &env.duplThreshold) {
		return false
	}
	if !getEnvInt("JPLAG_THRESHOLD", &env.jplagThreshold) {
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
