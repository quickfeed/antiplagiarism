package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/autograde/antiplagiarism/proto"
)

type apServer struct {
	env envVariables
}

func (s *apServer) CheckPlagiarism(ctx context.Context, req *pb.ApRequest) (*pb.ApResponse, error) {
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

	if !success {
		return &pb.ApResponse{Success: false, Err: "Check the server command prompt for the error."}, nil
	}

	return &pb.ApResponse{Success: true, Err: ""}, nil
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
	pb.RegisterApServer(grpcServer, server)
	fmt.Printf("Preparing to serve incoming requests.\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
