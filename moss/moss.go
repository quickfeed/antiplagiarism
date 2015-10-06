package moss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"../common"
	"github.com/PuerkitoBio/goquery"
)

type matches struct {
	url        string
	match1text string
	match2text string
}

// CreateCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. CreateCommands takes as input
// labsBaseDir, the location of the student directories,
// toolDir, the location of the Moss script, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling Moss
// to ignore matches that appear in at least that many files.
func CreateCommands(labsBaseDir string, toolDir string, labs []common.LabInfo, threshold int) ([]string, bool) {
	studentsLabDirs, success := directoryContents(labsBaseDir, labs)
	if !success {
		fmt.Printf("Error getting the student directories.\n")
		return nil, false
	}

	commands, success := createMossCommands(toolDir, studentsLabDirs, labs, threshold)
	if !success {
		fmt.Printf("Error creating the Moss commands.\n")
		return nil, false
	}

	return commands, true
}

// directoryContents returns a two-dimensional slice with full-path directories
// to send to Moss for evaluation. The first index addresses the specific lab,
// and the second index addresses the specific student. If a student does not
// have the lab directory, the 2D slice will save the directory as an empty
// string. The second return argument indicates whether or not the function was
// successful. directoryContents takes as input baseDir, the location of the
// student directories, and labs, a slice of the labs.
func directoryContents(baseDir string, labs []common.LabInfo) ([][]string, bool) {
	// Try to read the base directory
	contents, err := ioutil.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", baseDir, err)
		return nil, false
	}

	var studentDirs []string
	// Get a list of all the student directories (full path)
	for _, item := range contents {
		if item.IsDir() {
			studentDirs = append(studentDirs, filepath.Join(baseDir, item.Name()))
		}
	}

	studentsLabDirs := make([][]string, len(labs))
	// For each lab
	for i := range studentsLabDirs {
		studentsLabDirs[i] = make([]string, len(studentDirs))

		// For each student
		for j := range studentsLabDirs[i] {
			tempDir := filepath.Join(studentDirs[j], labs[i].Name)
			_, err := ioutil.ReadDir(tempDir)
			if err != nil {
				studentsLabDirs[i][j] = ""
			} else {
				studentsLabDirs[i][j] = tempDir
			}
		}
	}

	return studentsLabDirs, true
}

// createMossCommands will create Moss commands to upload the lab files.
// It returns a slice of Moss commands.  The second return argument indicates
// whether or not the function was successful. createMossCommands takes as input
// mossDir, the location of the Moss script, studentsLabDirs, a 2D slice of
// directories, labs, a slice of the labs, and threshold, an integer telling Moss
// to ignore matches that appear in at least that many files.
func createMossCommands(mossDir string, studentsLabDirs [][]string, labs []common.LabInfo, threshold int) ([]string, bool) {
	var commands []string
	mOption := "-m " + strconv.Itoa(threshold)

	// For each lab
	for i := range studentsLabDirs {
		var lOption string
		var fileExt []string

		// Set language option and file extensions
		if labs[i].Language == common.Golang {
			lOption = "-l java"
			fileExt = append(fileExt, "*.go")
		} else if labs[i].Language == common.Cpp {
			lOption = "-l cc"
			fileExt = append(fileExt, "*.cpp")
			fileExt = append(fileExt, "*.h")
		} else {
			lOption = "-l java"
			fileExt = append(fileExt, "*.java")
		}

		// Start creating the moss command
		var buf bytes.Buffer
		buf.WriteString(filepath.Join(mossDir, "moss") + " " + lOption + " " + mOption + " -d")

		// For each student
		for j := range studentsLabDirs[i] {

			// If student has the lab
			if studentsLabDirs[i][j] != "" {

				// Add all the files with the appropriate extensions
				for k := range fileExt {
					buf.WriteString(" " + filepath.Join(studentsLabDirs[i][j], fileExt[k]))
				}
			}
		}

		buf.WriteString(" > " + labs[i].Name + ".txt &")

		// Add the Moss command for this lab
		commands = append(commands, buf.String())
	}

	return commands, true
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
		linkBody, success := GetHTMLData(match.url)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+".html"), []byte(linkBody), 0644)

		// Get and save top frame
		topBody, success := GetHTMLData(topFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-top.html"), []byte(topBody), 0644)

		// Get and save left frame
		leftBody, success := GetHTMLData(leftFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-0.html"), []byte(leftBody), 0644)

		// Get and save right frame
		rightBody, success := GetHTMLData(rightFrame)
		if !success {
			return false
		}
		ioutil.WriteFile(filepath.Join(resultsDir, base+"-1.html"), []byte(rightBody), 0644)
	}

	MakeResultsMainPage(resultsDir, lab, comparisons)

	return true
}

// GetHTMLData returns html data as a string. It returns the html
// and whether or not the function was successful. GetHTMLData takes
// as input url, the location of the html.
func GetHTMLData(url string) (string, bool) {
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

// MakeResultsMainPage creates the main html file for the results.
// It takes as input resultsDir, where to save the data,
// lab, information about the current lab, and comparisons,
// information about the matches found.
func MakeResultsMainPage(resultsDir string, lab common.LabInfo, comparisons []matches) {
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
