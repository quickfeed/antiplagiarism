package moss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../common"
	"github.com/PuerkitoBio/goquery"
)

type matches struct {
	url        string
	match1text string
	match2text string
}

// SaveMossResults saves the data from the specified Moss URL. It returns
// whether or not the function was successful. SaveMossResults takes as input
// url, the main url for the Moss results, baseDir, where to save the data,
// and lab, information about the current lab.
func SaveMossResults(url string, baseDir string, lab common.LabInfo) bool {
	resultsDir := filepath.Join(baseDir, lab.Name)

	os.MkdirAll(resultsDir, 0764)
	os.Remove(filepath.Join(resultsDir, "*.*"))

	var comparisons []matches

	// Get web page data
	doc, err := goquery.NewDocument(url)
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

	makeResultsMainPage(resultsDir, lab, comparisons)

	return true
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
// lab, information about the current lab, and comparisons,
// information about the matches found.
func makeResultsMainPage(resultsDir string, lab common.LabInfo, comparisons []matches) {
	var buf bytes.Buffer
	buf.WriteString("<HTML>\n<HEAD>\n<TITLE>")
	buf.WriteString(lab.Name + " Results")
	buf.WriteString("</TITLE>\n</HEAD>\n<BODY>\n")
	buf.WriteString(lab.Name + " Results<br>")
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

	ioutil.WriteFile(filepath.Join(resultsDir, "results.html"), []byte(buf.String()), 0644)
}
