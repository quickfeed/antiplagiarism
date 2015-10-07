package common

// The language the lab was written in
const (
	Java = iota
	Golang
	Cpp
)

// LabInfo is information about the lab: name and programming language
type LabInfo struct {
	Name     string
	Class    string
	Language int
}
