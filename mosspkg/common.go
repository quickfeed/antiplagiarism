package mosspkg

const (
	Java = iota
	Golang
	Cpp
)

type LabInfo struct {
	Name     string
	Language int
}
