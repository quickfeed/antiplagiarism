package jplag

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

	// Find the last table
	table := doc.Find("table").Last()

	// Find the table rows
	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		repo1 := ""

		// Find the table columns
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			url := ""
			value := td.Text()

			// Find the href attributes
			td.Find("a").Each(func(k int, a *goquery.Selection) {
				url, _ = a.Attr("href")
			})

			// If there is a url to the match
			if url != "" {
				url = filepath.Join(resultsDir, url)

				tmpValue := strings.TrimSuffix(value, "%)")
				parts := strings.Split(tmpValue, "(")

				// Get the second student's repo name
				repo2 := parts[0]

				// Get the percentage
				pct, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					fmt.Printf("Error parsing %s: %s\n", parts[1], err)
				}

				// Save the percentage and url if the percentage if greater than the previous value.
				if pct > results[repo1].Percent {
					results[repo1] = common.ResultEntry{Repo: repo1, Percent: pct, URL: url}
				}
				if pct > results[repo2].Percent {
					results[repo2] = common.ResultEntry{Repo: repo2, Percent: pct, URL: url}
				}

			} else if value != "->" {
				// Get the first student's repo name
				repo1 = value
			}
		})
	})

	var orderedResults common.ResultEntries
	common.OrderResults(&results, &orderedResults)
	common.MakeResultsFile(resultsDir, &orderedResults)

	return true
}
