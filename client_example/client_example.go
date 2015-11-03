package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/autograde/antiplagiarism/proto"
)

// The language the lab was written in
const (
	Java = iota
	Golang
	Cpp
	C
)

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	endpoint = flag.String(
		"endpoint",
		"localhost:12111",
		"Endpoint on which server runs or to which client connects",
	)
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	// Create connection
	conn, err := grpc.Dial(*endpoint)
	if err != nil {
		fmt.Printf("Error while connecting to server: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to server on %v\n", *endpoint)

	// Create client
	client := pb.NewApClient(conn)

	// Create request
	request := pb.ApRequest{GithubOrg: "test-repo",
		GithubToken:  "12345",
		StudentRepos: []string{"student1-labs", "student2-labs"},
		LabNames:     []string{"lab1", "lab2"},
		LabLanguages: []int32{C, Golang}}

	// Send request and get response
	response, err := client.CheckPlagiarism(context.Background(), &request)

	// Check response
	if err != nil {
		fmt.Printf("gRPC error: %s\n", err)
	} else if response.Success == false {
		fmt.Printf("Anti-plagiarism error: %s\n", response.Err)
	} else {
		fmt.Printf("Anti-plagiarism application ran successfully.\n")
	}
}
