package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/wshaman/course-grpc/cmd/server/db"
	"github.com/wshaman/course-grpc/common/transport"
)

const (
	port = ":50051"
)

// server is used to implement transport.UserSearcher.
type server struct {
	transport.UnimplementedUserRegisterServer
}

// Search implements transport.UserSearcher
func (s *server) UserRegister(ctx context.Context, in *transport.UserRegisterRequest) (*transport.UserRegisterResponse, error) {
	log.Printf("Received: %v", in)
	if db.CheckIfExists(in) {
		log.Printf("User already in database")
		return nil, fmt.Errorf("User already exist ")
	}
	db.AddToDB(in)
	response := &transport.UserRegisterResponse{
		Result: db.Data.Map[in.GetName()].GetId(),
	}
	return response, nil
}

// List implements transport.List
func (s *server) UsersList(ctx context.Context, in *transport.EmptyRequest) (*transport.UsersListResponse, error) {
	p := db.LoadPersistent()
	result := &transport.UsersListResponse{
		Result: p,
	}
	return result, nil
}

func main() {
	db.Populate()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	transport.RegisterUserRegisterServer(s, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
