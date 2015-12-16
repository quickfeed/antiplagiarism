package jplag

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/autograde/antiplagiarism/common"
)

type jplagInputTest struct {
	org      string
	labs     []common.LabInfo
	commands []string
	success  bool
}

func TestJplagInput(t *testing.T) {
	dir, _ := filepath.Abs("../")
	threshStr := os.Getenv("JPLAG_THRESHOLD")
	threshold, err := strconv.Atoi(threshStr)
	if err != nil {
		t.Errorf("Set the JPLAG_THRESHOLD environment variable before running the TestJplagInput test.")
		return
	}
	toolLoc := os.Getenv("JPLAG_FULLY_QUALIFIED_NAME")
	if toolLoc == "" {
		t.Errorf("Set the JPLAG_FULLY_QUALIFIED_NAME environment variable before running the TestJplagInput test.")
		return
	}

	jplagInputTests := []jplagInputTest{
		// One lab, has java code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0}},
			[]string{"java -jar " + toolLoc + " -t " + threshStr + " -l java17 -r JPLAG.testOrg.lab1 -S lab1 -s " + filepath.Join(dir, "testOrg")},
			true},
		// One lab, has go code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab2", 1}},
			[]string{""},
			true},
		// One lab, has c++ code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab3", 2}},
			[]string{"java -jar " + toolLoc + " -t " + threshStr + " -l c/c++ -r JPLAG.testOrg.lab3 -S lab3 -s " + filepath.Join(dir, "testOrg")},
			true},
		// One lab, has c code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab4", 3}},
			[]string{"java -jar " + toolLoc + " -t " + threshStr + " -l c/c++ -r JPLAG.testOrg.lab4 -S lab4 -s " + filepath.Join(dir, "testOrg")},
			true},
		// Four labs, java, go, c++, c code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0},
				common.LabInfo{"lab2", 1},
				common.LabInfo{"lab3", 2},
				common.LabInfo{"lab4", 3}},
			[]string{
				"java -jar " + toolLoc + " -t " + threshStr + " -l java17 -r JPLAG.testOrg.lab1 -S lab1 -s " + filepath.Join(dir, "testOrg"),
				"",
				"java -jar " + toolLoc + " -t " + threshStr + " -l c/c++ -r JPLAG.testOrg.lab3 -S lab3 -s " + filepath.Join(dir, "testOrg"),
				"java -jar " + toolLoc + " -t " + threshStr + " -l c/c++ -r JPLAG.testOrg.lab4 -S lab4 -s " + filepath.Join(dir, "testOrg")},
			true},
		// Directory doesn't exist
		{"DAT500",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0},
				common.LabInfo{"lab2", 1},
				common.LabInfo{"lab3", 2},
				common.LabInfo{"lab4", 3}},
			nil,
			false},
	}

	jplagTool := Jplag{LabsBaseDir: dir, ToolFqn: toolLoc, Threshold: threshold}

	for i, tst := range jplagInputTests {
		commands, success := jplagTool.CreateCommands(tst.org, tst.labs)
		sort.Strings(commands)
		sort.Strings(tst.commands)
		if !common.CompareStringSlices(commands, tst.commands) || success != tst.success {
			t.Errorf("Jplag Input Test %d: \n\tinput %s, %v\n\tgot: %v, %v\n\twnt: %v, %v", i, tst.org, tst.labs, success, commands, tst.success, tst.commands)
		}
	}
}
