package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	apProto "github.com/autograde/antiplagiarism/proto"
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
		"localhost:11111",
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

	var opts []grpc.DialOption
	// Note that this is not secure.
	opts = append(opts, grpc.WithInsecure())

	// Create connection
	conn, err := grpc.Dial(*endpoint, opts...)
	if err != nil {
		fmt.Printf("Error while connecting to server: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to server on %v\n", *endpoint)

	// Create client
	client := apProto.NewApClient(conn)

	labs := []*apProto.ApRequestLab{
		&apProto.ApRequestLab{Name: "lab1", Language: int32(C)},
		&apProto.ApRequestLab{Name: "lab2", Language: int32(Golang)}}

	// Create request
	request := apProto.ApRequest{GithubOrg: "test-repo",
		GithubToken:  "testtoken",
		StudentRepos: []string{"student1-labs", "student2-labs"},
		Labs:         labs}

	// Send request and get response
	response, err := client.CheckPlagiarism(context.Background(), &request)

	// Check response
	if err != nil {
		fmt.Printf("gRPC error: %s\n", err)
	} else if response.Success == false {
		fmt.Printf("Anti-plagiarism error: %s\n", response.Err)
	} else {
		fmt.Printf("Anti-plagiarism application received request.\n")
	}
}
