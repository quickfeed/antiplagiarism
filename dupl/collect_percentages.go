package dupl

import (
	"bufio"
	"fmt"
	"os"
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

	doc.Find("h2").Each(func(i int, h2 *goquery.Selection) {
		value := h2.Text()
		pos1 := strings.Index(value, org+"/")
		repo := value[pos1+len(org)+1:]
		pos2 := strings.Index(repo, "/")
		repo = repo[:pos2]
		pct := 1.0
		url := fileNameAndPath

		results[repo] = common.ResultEntry{Repo: repo, Percent: pct, URL: url}
	})

	var orderedResults common.ResultEntries
	common.OrderResults(&results, &orderedResults)
	common.MakeResultsFile(resultsDir, &orderedResults)

	return true
}
