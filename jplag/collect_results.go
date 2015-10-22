package jplag

import (
	"../common"
)

// SaveResults saves the data from the specified file. It returns
// whether or not the function was successful. SaveResults takes as input
// fileNameAndPath, the file output by dupl, baseDir, where to save the data,
// and lab, information about the current lab.
func (j Jplag) SaveResults(fileNameAndPath string, baseDir string, lab common.LabInfo) bool {
	return true
}
