package moss

import (
	"../common"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// collectPercentages searches the output file for the highest
// plagiarism percentages for each student.
func collectPercentages(org string, fileNameAndPath string) bool {
	results := make(map[string]common.PctInfo)
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

		// If there is a link to the match
		if url != "" {
			value := a.Text()
			data := strings.Split(value, ", ")

			for i := range data {
				temp, _ := getInfoFromString(org, data[i])
				temp.Link = url
				idx := temp.Repo

				// Save the percentage and link if the percentage if greater than the previous value.
				if temp.Percent > results[idx].Percent {
					results[idx] = temp
				}
			}
		}
	})

	var orderedResults common.PctInfoSlice
	common.OrderPctInfo(&results, &orderedResults)

	common.MakePercentagePage(resultsDir, &orderedResults)
	common.MakePercentageFile(resultsDir, &orderedResults)

	return true
}

// getInfoFromString gets the repo name and percentage from a string.
// It returns a PctInfo containing the information (without the link) and
// whether or not the function was successful. It takes as input org,
// the name of the GitHub organization and rawData, the string.
func getInfoFromString(org string, rawData string) (common.PctInfo, bool) {

	// Split the data into location and percentage
	dataParts := strings.Split(rawData, " (")

	// Get the repo from the location
	dataParts[1] = strings.TrimSpace(dataParts[1])
	dataParts[1] = strings.TrimSuffix(dataParts[1], "%)")

	// Get the percentage
	pct, err := strconv.ParseFloat(dataParts[1], 64)
	if err != nil {
		fmt.Printf("Error parsing %s: %s\n", dataParts[1], err)
		return common.PctInfo{}, false
	}

	pos1 := strings.Index(dataParts[0], org+"/")
	repo := dataParts[0][pos1+len(org)+1:]
	pos2 := strings.Index(repo, "/")
	repo = repo[:pos2]

	return common.PctInfo{Repo: repo, Percent: pct, Link: ""}, true
}
