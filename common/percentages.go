package common

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
)

// ResultEntries is a slice of ResultEntry
type ResultEntries []ResultEntry

// ResultEntry contains plagiarism information about a repository
type ResultEntry struct {
	Repo    string
	Percent float64
	Link    string
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

// MakePercentagePage saves the highest percentage results for each student
// and a link to code in an html file called percentage.html.
// It takes as input resultsDir, the location to save the file,
// and results, a pointer to the slice of percentage information
func MakePercentagePage(resultsDir string, results *ResultEntries) {
	var buf bytes.Buffer
	buf.WriteString("<HTML>\n<HEAD>\n<TITLE>")
	buf.WriteString("Percentages")
	buf.WriteString("</TITLE>\n</HEAD>\n<BODY>\n")
	buf.WriteString("Percentages<br><br>")
	for _, result := range *results {
		buf.WriteString("\n<A HREF=\"")
		buf.WriteString(result.Link)
		buf.WriteString("\">")
		buf.WriteString(result.Repo + " " + strconv.FormatFloat(result.Percent, 'f', 1, 64))
		buf.WriteString("</A><br>")
	}
	buf.WriteString("\n</BODY>\n</HTML>\n")

	ioutil.WriteFile(filepath.Join(resultsDir, "percentage.html"), []byte(buf.String()), 0644)
}

// MakePercentageFile saves the highest percentage results for each student
// and a link to code in a text file called percentage.txt
// It takes as input resultsDir, the location to save the file,
// and results, a pointer to the slice of percentage information
func MakePercentageFile(resultsDir string, results *ResultEntries) {
	var buf bytes.Buffer
	for _, result := range *results {
		buf.WriteString(result.Repo)
		buf.WriteString("|")
		buf.WriteString(strconv.FormatFloat(result.Percent, 'f', 1, 64))
		buf.WriteString("|")
		buf.WriteString(result.Link)
		buf.WriteString("\n")
	}

	ioutil.WriteFile(filepath.Join(resultsDir, "percentage.txt"), []byte(buf.String()), 0644)
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
