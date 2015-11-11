package moss

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/autograde/antiplagiarism/common"
)

type matches struct {
	url        string
	match1text string
	match2text string
}

// SaveResults looks for Moss results and saves them. It returns
// whether or not the function was successful. SaveResults takes as input
// org, the GitHub organization name, and labs, a slice of the labs,
// path, where Moss puts the results, and baseDir, where to save the data.
func (m Moss) SaveResults(org string, labs []common.LabInfo, path string, baseDir string) bool {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", path, err)
		return false
	}

	// For each lab
	for _, lab := range labs {
		fileName := "MOSS." + org + "." + lab.Name + ".txt"
		fileNameAndPath := filepath.Join(path, fileName)
		success := saveLabResults(fileNameAndPath, baseDir, org, lab.Name)

		if success {
			os.Remove(fileNameAndPath)
		}
	}

	return true
}

// saveLabResults saves the data from the specified Moss URL. It returns
// whether or not the function was successful. SaveResults takes as input
// fileNameAndPath, where Moss saved the results of this lab,
// baseDir, where to save the data, org, the name of the organization (class)
// and labName, the name of this lab.
func saveLabResults(fileNameAndPath string, baseDir string, org string, labName string) bool {
	resultsDir := filepath.Join(baseDir, org, labName, "moss")

	// Remove old results.
	err := os.RemoveAll(resultsDir)
	if err != nil {
		fmt.Printf("Error removing the old results directory. %s\n", err)
		return false
	}

	// Make the directory.
	err = os.MkdirAll(resultsDir, 0764)
	if err != nil {
		fmt.Printf("Error creating the results directory. %s\n", err)
		return false
	}

	resultsURL, success := getURLFromOutput(fileNameAndPath)
	if !success {
		return false
	}

	var comparisons []matches

	// Get web page data
	doc, err := goquery.NewDocument(resultsURL)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}

	// Find the table rows
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var match matches

		// Find the table columns
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			url := ""

			// Find the href attributes
			td.Find("a").Each(func(k int, a *goquery.Selection) {
				url, _ = a.Attr("href")
			})

			// If there is an href attribute
			if url != "" {
				match.url = url
				val := td.Text()
				if match.match1text == "" {
					match.match1text = val
				} else {
					match.match2text = val
				}
			}
		})

		if match.url != "" {
			comparisons = append(comparisons, match)
		}
	})

	// For each URL
	for _, match := range comparisons {
		pos1 := strings.LastIndex(match.url, "/")
		pos2 := strings.LastIndex(match.url, ".html")
		base := match.url[pos1+1 : pos2]
		topFrame := strings.Replace(match.url, ".html", "-top.html", 1)
		leftFrame := strings.Replace(match.url, ".html", "-0.html", 1)
		rightFrame := strings.Replace(match.url, ".html", "-1.html", 1)

		// Get and save web page data
		linkBody, success := getHTMLData(match.url)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+".html"), []byte(linkBody), 0644)

		// Get and save top frame
		topBody, success := getHTMLData(topFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-top.html"), []byte(topBody), 0644)

		// Get and save left frame
		leftBody, success := getHTMLData(leftFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-0.html"), []byte(leftBody), 0644)

		// Get and save right frame
		rightBody, success := getHTMLData(rightFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-1.html"), []byte(rightBody), 0644)
	}

	makeResultsMainPage(resultsDir, labName, comparisons)
	collectPercentages(org, filepath.Join(resultsDir, "index.html"))

	return true
}

// getURLFromOutput returns the url from the Moss output. It returns the url
// and whether or not the function was successful. getURLFromOutput takes
// as input fileNameAndPath, the location of the output.
func getURLFromOutput(fileNameAndPath string) (string, bool) {
	foundURL := false
	url := ""

	file, err := os.Open(fileNameAndPath)
	if err != nil {
		fmt.Printf("Error reading output from Moss. %s\n", err)
		return url, foundURL
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "http://moss.stanford.edu/results/") {
			url = line
			foundURL = true
			break
		}
	}

	return url, foundURL
}

// getHTMLData returns html data as a string. It returns the html
// and whether or not the function was successful. getHTMLData takes
// as input url, the location of the html.
func getHTMLData(url string) (string, bool) {
	var body string
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("%v\n", err)
		return "", false
	} else if resp.StatusCode != 200 {
		fmt.Printf("%v\n", resp.Status)
		return "", false
	} else {
		bodyBuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v\n", err)
			return "", false
		}

		body = string(bodyBuf)
	}

	return body, true
}

// makeResultsMainPage creates the main html file for the results.
// It takes as input resultsDir, where to save the data,
// labName, the name of the current lab, and comparisons,
// information about the matches found.
func makeResultsMainPage(resultsDir string, labName string, comparisons []matches) {
	var buf bytes.Buffer
	buf.WriteString("<HTML>\n<HEAD>\n<TITLE>")
	buf.WriteString(labName + " Results")
	buf.WriteString("</TITLE>\n</HEAD>\n<BODY>\n")
	buf.WriteString(labName + " Results<br><br>")
	for _, match := range comparisons {
		pos := strings.LastIndex(match.url, "/")
		base := match.url[pos+1:]

		buf.WriteString("\n<A HREF=\"")
		buf.WriteString(base)
		buf.WriteString("\">")
		buf.WriteString(match.match1text + ", " + match.match2text)
		buf.WriteString("</A><br>")
	}
	buf.WriteString("\n</BODY>\n</HTML>\n")

	ioutil.WriteFile(filepath.Join(resultsDir, "index.html"), []byte(buf.String()), 0644)
}
