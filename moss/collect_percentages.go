package moss

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/autograde/antiplagiarism/common"
)

// collectPercentages searches the output file for the highest
// plagiarism percentages for each student.
func collectPercentages(org string, fileNameAndPath string) bool {
	results := make(map[string]common.ResultEntry)
	pos := strings.LastIndex(fileNameAndPath, "/")
	resultsDir := fileNameAndPath[:pos]

	file, err := os.Open(fileNameAndPath)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", fileNameAndPath, err)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Get web page data
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Printf("Error opening document in goquery: %v\n", err)
		return false
	}

	// Find the href attributes
	doc.Find("a").Each(func(k int, a *goquery.Selection) {
		url, _ := a.Attr("href")

		// If there is a url to the match
		if url != "" {
			value := a.Text()
			data := strings.Split(value, ", ")

			for i := range data {
				temp, _ := getInfoFromString(org, data[i])
				temp.URL = filepath.Join(resultsDir, url)
				idx := temp.Repo

				// Save the percentage and url if the percentage if greater than the previous value.
				if temp.Percent > results[idx].Percent {
					results[idx] = temp
				}
			}
		}
	})

	var orderedResults common.ResultEntries
	common.OrderResults(&results, &orderedResults)
	common.MakeResultsFile(resultsDir, &orderedResults)

	return true
}

// getInfoFromString gets the repo name and percentage from a string.
// It returns a ResultEntry containing the information (without the url) and
// whether or not the function was successful. It takes as input org,
// the name of the GitHub organization and rawData, the string.
func getInfoFromString(org string, rawData string) (common.ResultEntry, bool) {

	// Split the data into location and percentage
	dataParts := strings.Split(rawData, " (")

	// Get the repo from the location
	dataParts[1] = strings.TrimSpace(dataParts[1])
	dataParts[1] = strings.TrimSuffix(dataParts[1], "%)")

	// Get the percentage
	pct, err := strconv.ParseFloat(dataParts[1], 64)
	if err != nil {
		fmt.Printf("Error parsing %s: %s\n", dataParts[1], err)
		return common.ResultEntry{}, false
	}

	pos1 := strings.Index(dataParts[0], org+"/")
	repo := dataParts[0][pos1+len(org)+1:]
	pos2 := strings.Index(repo, "/")
	repo = repo[:pos2]

	return common.ResultEntry{Repo: repo, Percent: pct, URL: ""}, true
}
