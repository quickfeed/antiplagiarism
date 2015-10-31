package common

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
)

// PctInfoSlice is a slice of PctInfo
type PctInfoSlice []PctInfo

// PctInfo contains plagiarism information about a repository
type PctInfo struct {
	Repo    string
	Percent float64
	Link    string
}

// Len is the number of elements in the collection.
func (p PctInfoSlice) Len() int {
	return len(p)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (p PctInfoSlice) Less(i, j int) bool {
	return p[i].Repo < p[j].Repo
}

// Swap swaps the elements with indexes i and j.
func (p PctInfoSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// MakePercentagePage saves the highest percentage results for each student
// and a link to code in an html file called percentage.html.
// It takes as input resultsDir, the location to save the file,
// and results, a pointer to the slice of percentage information
func MakePercentagePage(resultsDir string, results *PctInfoSlice) {
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
func MakePercentageFile(resultsDir string, results *PctInfoSlice) {
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

// OrderPctInfo converts a map of PctInfo into a sorted slice of PctInfo.
// It takes as input input, a map of PctInfo, and output, a slice of PctInfo.
func OrderPctInfo(input *map[string]PctInfo, output *PctInfoSlice) {
	*output = nil
	for _, result := range *input {
		*output = append(*output, result)
	}

	sort.Sort(*output)
}
