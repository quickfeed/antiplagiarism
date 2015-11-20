package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

// ResultsFileName is the name of the results file
var ResultsFileName = "results.json"

// ResultEntries is a slice of ResultEntry
type ResultEntries []ResultEntry

// ResultEntry contains plagiarism information about a repository
type ResultEntry struct {
	Repo    string
	Percent float64
	URL     string
}

// Len is the number of elements in the collection.
func (r ResultEntries) Len() int {
	return len(r)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (r ResultEntries) Less(i, j int) bool {
	return r[i].Repo < r[j].Repo
}

// Swap swaps the elements with indexes i and j.
func (r ResultEntries) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// MakeResultsFile saves the highest percentage results for each student
// and a link to code in a text file called percentage.txt
// It takes as input resultsDir, the location to save the file,
// and results, a pointer to the slice of results
func MakeResultsFile(resultsDir string, results *ResultEntries) {
	// Encode Message into JSON format
	buf, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("Error marshalling results in JSON format. %s\n", err)
		return
	}

	ioutil.WriteFile(filepath.Join(resultsDir, ResultsFileName), buf, 0644)
}

// OrderResults converts a map of ResultEntry into a sorted slice of ResultEntry.
// It takes as input input, a map of ResultEntry, and output, a slice of ResultEntry.
func OrderResults(input *map[string]ResultEntry, output *ResultEntries) {
	*output = nil
	for _, result := range *input {
		*output = append(*output, result)
	}

	sort.Sort(*output)
}
