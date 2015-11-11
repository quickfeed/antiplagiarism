package dupl

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/autograde/antiplagiarism/common"
)

type duplInputTest struct {
	org      string
	labs     []common.LabInfo
	commands []string
	success  bool
}

func TestDuplInput(t *testing.T) {
	dir, _ := filepath.Abs("../")
	threshStr := os.Getenv("DUPL_THRESHOLD")
	threshold, err := strconv.Atoi(threshStr)
	if err != nil {
		t.Errorf("Set the DUPL_THRESHOLD environment variable before running the TestDuplInput test.")
		return
	}

	duplInputTests := []duplInputTest{
		// One lab, has go code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab2", 1}},
			[]string{"dupl -t " + threshStr + " -html " +
				dir + "/testOrg/student1/lab2/ " +
				dir + "/testOrg/student2/lab2/ " +
				dir + "/testOrg/student3/lab2/ " +
				dir + "/testOrg/student4/lab2/ " +
				dir + "/testOrg/student5/lab2/ > DUPL.testOrg.lab2.html"},
			true},
		// One lab, doesn't have go code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab3", 2}},
			[]string{""},
			true},
		// Four labs, only lab 2 has go code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0},
				common.LabInfo{"lab2", 1},
				common.LabInfo{"lab3", 2},
				common.LabInfo{"lab4", 3}},
			[]string{"",
				"dupl -t " + threshStr + " -html " +
					dir + "/testOrg/student1/lab2/ " +
					dir + "/testOrg/student2/lab2/ " +
					dir + "/testOrg/student3/lab2/ " +
					dir + "/testOrg/student4/lab2/ " +
					dir + "/testOrg/student5/lab2/ > DUPL.testOrg.lab2.html",
				"",
				""},
			true},
		// Four labs, all four have go code
		/*{"testOrg",
		[]common.LabInfo{
			common.LabInfo{"lab1", 1},
			common.LabInfo{"lab2", 1},
			common.LabInfo{"lab3", 1},
			common.LabInfo{"lab4", 1}},
		[]string{"dupl -t " + threshStr + " -html " +
		dir + "/testOrg/student1/lab1/ " +
		dir + "/testOrg/student2/lab1/ " +
		dir + "/testOrg/student3/lab1/ " +
		dir + "/testOrg/student4/lab1/ " +
		dir + "/testOrg/student5/lab1/ > DUPL.testOrg.lab1.html",
			"dupl -t " + threshStr + " -html " +
			dir + "/testOrg/student1/lab2/ " +
			dir + "/testOrg/student2/lab2/ " +
			dir + "/testOrg/student3/lab2/ " +
			dir + "/testOrg/student4/lab2/ " +
			dir + "/testOrg/student5/lab2/ > DUPL.testOrg.lab2.html",
			"dupl -t " + threshStr + " -html " +
			dir + "/testOrg/student1/lab3/ " +
			dir + "/testOrg/student2/lab3/ " +
			dir + "/testOrg/student3/lab3/ " +
			dir + "/testOrg/student4/lab3/ " +
			dir + "/testOrg/student5/lab3/ > DUPL.testOrg.lab2.html",
			"dupl -t " + threshStr + " -html " +
			dir + "/testOrg/student1/lab4/ " +
			dir + "/testOrg/student2/lab4/ " +
			dir + "/testOrg/student3/lab4/ " +
			dir + "/testOrg/student4/lab4/ " +
			dir + "/testOrg/student5/lab4/ > DUPL.testOrg.lab3.html"},
		true},*/
		// Four labs, no labs have go code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0},
				common.LabInfo{"lab2", 2},
				common.LabInfo{"lab3", 2},
				common.LabInfo{"lab4", 3}},
			[]string{"",
				"",
				"",
				""},
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

	duplTool := Dupl{LabsBaseDir: dir, ToolFqn: "", Threshold: threshold}

	for i, tst := range duplInputTests {
		commands, success := duplTool.CreateCommands(tst.org, tst.labs)
		sort.Strings(commands)
		sort.Strings(tst.commands)
		if !common.CompareStringSlices(commands, tst.commands) || success != tst.success {
			t.Errorf("Dupl Input Test %d: \n\tinput %s, %v\n\tgot: %v, %v\n\twnt: %v, %v", i, tst.org, tst.labs, success, commands, tst.success, tst.commands)
		}
	}
}
