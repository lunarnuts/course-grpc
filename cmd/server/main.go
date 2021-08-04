package main

import (
	"context"
	"fmt"
	"log"
	"net"

	cs "github.com/wshaman/contacts-stub"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/wshaman/course-grpc/common/transport"
)

const (
	port = ":50051"
)

// server is used to implement transport.UserSearcher.
type server struct {
	mp map[string]transport.Person
	db []transport.Person
	transport.UnimplementedUserRegisterServer
	transport.UnimplementedUserListServer
}

func (s server) LoadPersistent() []transport.Person {
	return s.db
}

func (s server) FindByName(person transport.RegisterRequest) (transport.Person, error) {
	val, ok := s.mp[person.Name]
	if !ok {
		return transport.Person{}, fmt.Errorf("not found")
	}
	return val, nil
}

func (s *server) AddToDB(person transport.RegisterRequest) error {
	if _, ok := s.FindByName(person); ok != nil {
		return ok
	}
	newPerson := transport.Person{
		Id:   int64(len(s.db)),
		Name: person.Name,
	}
	s.db = append(s.db, newPerson)
	s.mp[person.GetName()] = newPerson
	return nil
}

func (s *server) Populate() {
	var list = []transport.RegisterRequest{
		{Name: "Egor Uliyanov"},
		{Name: "ilya Komarskih"},
		{Name: "Ivan Ivanov"},
		{Name: "Alex Fergusson"},
		{Name: "Julian Casablancas"},
		{Name: "Pavel Lyadov"},
	}
	for _, person := range list {
		s.AddToDB(person)
	}
}

// Search implements transport.UserSearcher
func (s *server) Register(ctx context.Context, in *transport.RegisterRequest) (*transport.RegisterResponse, error) {
	log.Printf("Received: %v", in)
	res, err := s.FindByName(*in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find ")
	}
	response := &transport.RegisterResponse{
		Result: 
	}
	for _, v := range res {
		response.Result = append(response.Result, &transport.Person{
			Id:    int64(v.ID),
			Name:  v.LastName + " " + v.FirstName,
			Phone: v.Phone,
		})
	}
	return response, nil
}

func main() {
	p := cs.LoadPersistent()
	if err := cs.Populate(p); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	transport.RegisterUserRegisterServer(s, &server{})
	transport.RegisterUserListServer(s, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
