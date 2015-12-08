package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	apProto "github.com/autograde/antiplagiarism/proto"
)

type apServer struct {
	env envVariables
}

func (s *apServer) CheckPlagiarism(ctx context.Context, req *apProto.ApRequest) (*apProto.ApResponse, error) {
	fmt.Printf("Received request.\n")
	// If request is coming from the client example, just return.
	if req.GithubToken == "testtoken" {
		fmt.Printf("Sending response.\n")
		return &apProto.ApResponse{Success: true, Err: ""}, nil
	}

	var names []string
	var languages []int
	for i := range req.Labs {
		names = append(names, req.Labs[i].Name)
		languages = append(languages, int(req.Labs[i].Language))
	}

	args := commandLineArgs{
		studentRepos: req.StudentRepos,
		labNames:     names,
		labLanguages: languages,
		githubOrg:    req.GithubOrg,
		githubToken:  req.GithubToken,
		endpoint:     "",
	}

	success := buildAndRunCommands(&args, &s.env)
	fmt.Printf("Sending response.\n")

	if !success {
		return &apProto.ApResponse{Success: false, Err: "Check the server command prompt for the error."}, nil
	}

	return &apProto.ApResponse{Success: true, Err: ""}, nil
}

func startServer(args *commandLineArgs, env *envVariables) {
	listener, err := net.Listen("tcp", *endpoint)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Listener started on %v\n", *endpoint)
	}

	server := new(apServer)
	server.env = *env
	// TODO: Add transport security
	grpcServer := grpc.NewServer()
	apProto.RegisterApServer(grpcServer, server)
	fmt.Printf("Preparing to serve incoming requests.\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
