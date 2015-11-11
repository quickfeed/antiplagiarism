package moss

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/autograde/antiplagiarism/common"
)

type mossInputTest struct {
	org      string
	labs     []common.LabInfo
	commands [][]string
	success  bool
}

func TestMossInput(t *testing.T) {
	dir, _ := filepath.Abs("../")
	threshStr := os.Getenv("MOSS_THRESHOLD")
	threshold, err := strconv.Atoi(threshStr)
	if err != nil {
		t.Errorf("Set the MOSS_THRESHOLD environment variable before running the TestMossInput test.")
		return
	}
	toolLoc := os.Getenv("MOSS_FULLY_QUALIFIED_NAME")
	if toolLoc == "" {
		t.Errorf("Set the MOSS_FULLY_QUALIFIED_NAME environment variable before running the TestMossInput test.")
		return
	}

	mossInputTests := []mossInputTest{
		// One lab, has java code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0}},
			[][]string{[]string{toolLoc + " -l java -m " + threshStr + " -d ",
				dir + "/testOrg/student1/lab1/*.java ", dir + "/testOrg/student1/lab1/part2/*.java ",
				dir + "/testOrg/student2/lab1/*.java ", dir + "/testOrg/student2/lab1/part2/*.java ",
				dir + "/testOrg/student3/lab1/*.java ", dir + "/testOrg/student3/lab1/part2/*.java ",
				dir + "/testOrg/student4/lab1/*.java ", dir + "/testOrg/student4/lab1/part2/*.java ",
				dir + "/testOrg/student5/lab1/*.java ", dir + "/testOrg/student5/lab1/part2/*.java ",
				"> MOSS.testOrg.lab1.txt"}},
			true},
		// One lab, has go code, setting the language to Java is correct
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab2", 1}},
			[][]string{[]string{toolLoc + " -l java -m " + threshStr + " -d ",
				dir + "/testOrg/student1/lab2/*.go ", dir + "/testOrg/student1/lab2/part2/*.go ",
				dir + "/testOrg/student2/lab2/*.go ", dir + "/testOrg/student2/lab2/part2/*.go ",
				dir + "/testOrg/student3/lab2/*.go ", dir + "/testOrg/student3/lab2/part2/*.go ",
				dir + "/testOrg/student4/lab2/*.go ", dir + "/testOrg/student4/lab2/part2/*.go ",
				dir + "/testOrg/student5/lab2/*.go ", dir + "/testOrg/student5/lab2/part2/*.go ",
				"> MOSS.testOrg.lab2.txt"}},
			true},
		// One lab, has c++ code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab3", 2}},
			[][]string{[]string{toolLoc + " -l cc -m " + threshStr + " -d ",
				dir + "/testOrg/student1/lab3/*.cpp ", dir + "/testOrg/student1/lab3/*.h ",
				dir + "/testOrg/student1/lab3/part2/*.cpp ", dir + "/testOrg/student1/lab3/part2/*.h ",
				dir + "/testOrg/student2/lab3/*.cpp ", dir + "/testOrg/student2/lab3/*.h ",
				dir + "/testOrg/student2/lab3/part2/*.cpp ", dir + "/testOrg/student2/lab3/part2/*.h ",
				dir + "/testOrg/student3/lab3/*.cpp ", dir + "/testOrg/student3/lab3/*.h ",
				dir + "/testOrg/student3/lab3/part2/*.cpp ", dir + "/testOrg/student3/lab3/part2/*.h ",
				dir + "/testOrg/student4/lab3/*.cpp ", dir + "/testOrg/student4/lab3/*.h ",
				dir + "/testOrg/student4/lab3/part2/*.cpp ", dir + "/testOrg/student4/lab3/part2/*.h ",
				dir + "/testOrg/student5/lab3/*.cpp ", dir + "/testOrg/student5/lab3/*.h ",
				dir + "/testOrg/student5/lab3/part2/*.cpp ", // Student 5 does not have .h files in lab3/part2
				"> MOSS.testOrg.lab3.txt"}},
			true},
		// One lab, has c code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab4", 3}},
			[][]string{[]string{toolLoc + " -l c -m " + threshStr + " -d ",
				dir + "/testOrg/student1/lab4/*.c ", dir + "/testOrg/student1/lab4/*.h ",
				dir + "/testOrg/student1/lab4/part2/*.c ", dir + "/testOrg/student1/lab4/part2/*.h ",
				dir + "/testOrg/student2/lab4/*.c ", dir + "/testOrg/student2/lab4/*.h ",
				dir + "/testOrg/student2/lab4/part2/*.c ", dir + "/testOrg/student2/lab4/part2/*.h ",
				dir + "/testOrg/student3/lab4/*.c ", dir + "/testOrg/student3/lab4/*.h ",
				dir + "/testOrg/student3/lab4/part2/*.c ", dir + "/testOrg/student3/lab4/part2/*.h ",
				dir + "/testOrg/student4/lab4/*.c ", dir + "/testOrg/student4/lab4/*.h ",
				dir + "/testOrg/student4/lab4/part2/*.c ", dir + "/testOrg/student4/lab4/part2/*.h ",
				dir + "/testOrg/student5/lab4/*.c ", dir + "/testOrg/student5/lab4/*.h ",
				dir + "/testOrg/student5/lab4/part2/*.h ", // Student 5 does not have .c files in lab4/part2
				"> MOSS.testOrg.lab4.txt"}},
			true},
		// Four labs, java, go, c++, c code
		{"testOrg",
			[]common.LabInfo{
				common.LabInfo{"lab1", 0},
				common.LabInfo{"lab2", 1},
				common.LabInfo{"lab3", 2},
				common.LabInfo{"lab4", 3}},
			[][]string{
				[]string{toolLoc + " -l java -m " + threshStr + " -d ",
					dir + "/testOrg/student1/lab1/*.java ", dir + "/testOrg/student1/lab1/part2/*.java ",
					dir + "/testOrg/student2/lab1/*.java ", dir + "/testOrg/student2/lab1/part2/*.java ",
					dir + "/testOrg/student3/lab1/*.java ", dir + "/testOrg/student3/lab1/part2/*.java ",
					dir + "/testOrg/student4/lab1/*.java ", dir + "/testOrg/student4/lab1/part2/*.java ",
					dir + "/testOrg/student5/lab1/*.java ", dir + "/testOrg/student5/lab1/part2/*.java ",
					"> MOSS.testOrg.lab1.txt"},
				[]string{toolLoc + " -l java -m " + threshStr + " -d ",
					dir + "/testOrg/student1/lab2/*.go ", dir + "/testOrg/student1/lab2/part2/*.go ",
					dir + "/testOrg/student2/lab2/*.go ", dir + "/testOrg/student2/lab2/part2/*.go ",
					dir + "/testOrg/student3/lab2/*.go ", dir + "/testOrg/student3/lab2/part2/*.go ",
					dir + "/testOrg/student4/lab2/*.go ", dir + "/testOrg/student4/lab2/part2/*.go ",
					dir + "/testOrg/student5/lab2/*.go ", dir + "/testOrg/student5/lab2/part2/*.go ",
					"> MOSS.testOrg.lab2.txt"},
				[]string{toolLoc + " -l cc -m " + threshStr + " -d ",
					dir + "/testOrg/student1/lab3/*.cpp ", dir + "/testOrg/student1/lab3/*.h ",
					dir + "/testOrg/student1/lab3/part2/*.cpp ", dir + "/testOrg/student1/lab3/part2/*.h ",
					dir + "/testOrg/student2/lab3/*.cpp ", dir + "/testOrg/student2/lab3/*.h ",
					dir + "/testOrg/student2/lab3/part2/*.cpp ", dir + "/testOrg/student2/lab3/part2/*.h ",
					dir + "/testOrg/student3/lab3/*.cpp ", dir + "/testOrg/student3/lab3/*.h ",
					dir + "/testOrg/student3/lab3/part2/*.cpp ", dir + "/testOrg/student3/lab3/part2/*.h ",
					dir + "/testOrg/student4/lab3/*.cpp ", dir + "/testOrg/student4/lab3/*.h ",
					dir + "/testOrg/student4/lab3/part2/*.cpp ", dir + "/testOrg/student4/lab3/part2/*.h ",
					dir + "/testOrg/student5/lab3/*.cpp ", dir + "/testOrg/student5/lab3/*.h ",
					dir + "/testOrg/student5/lab3/part2/*.cpp ", // Student 5 does not have .h files in lab3/part2
					"> MOSS.testOrg.lab3.txt"},
				[]string{toolLoc + " -l c -m " + threshStr + " -d ",
					dir + "/testOrg/student1/lab4/*.c ", dir + "/testOrg/student1/lab4/*.h ",
					dir + "/testOrg/student1/lab4/part2/*.c ", dir + "/testOrg/student1/lab4/part2/*.h ",
					dir + "/testOrg/student2/lab4/*.c ", dir + "/testOrg/student2/lab4/*.h ",
					dir + "/testOrg/student2/lab4/part2/*.c ", dir + "/testOrg/student2/lab4/part2/*.h ",
					dir + "/testOrg/student3/lab4/*.c ", dir + "/testOrg/student3/lab4/*.h ",
					dir + "/testOrg/student3/lab4/part2/*.c ", dir + "/testOrg/student3/lab4/part2/*.h ",
					dir + "/testOrg/student4/lab4/*.c ", dir + "/testOrg/student4/lab4/*.h ",
					dir + "/testOrg/student4/lab4/part2/*.c ", dir + "/testOrg/student4/lab4/part2/*.h ",
					dir + "/testOrg/student5/lab4/*.c ", dir + "/testOrg/student5/lab4/*.h ",
					dir + "/testOrg/student5/lab4/part2/*.h ", // Student 5 does not have .c files in lab4/part2
					"> MOSS.testOrg.lab4.txt"}},
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

	mossTool := Moss{LabsBaseDir: dir, ToolFqn: toolLoc, Threshold: threshold}

	for i, tst := range mossInputTests {
		commands, success := mossTool.CreateCommands(tst.org, tst.labs)
		if !common.CompareStringSlicesRandom(commands, tst.commands) || success != tst.success {
			t.Errorf("Moss Input Test %d: \n\tinput %s, %v\n\tgot: %v, %v\n\twnt: %v, %v", i, tst.org, tst.labs, success, commands, tst.success, tst.commands)
		}
	}
}
