package main

import (
	"context"
	cs "github.com/wshaman/contacts-stub"
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/wshaman/course-grpc/common/transport"
)

const (
	port = ":50051"
)

// server is used to implement transport.UserSearcher.
type server struct {
	transport.UnimplementedUserSearcherServer
}

// Search implements transport.UserSearcher
func (s *server) Search(ctx context.Context, in *transport.SearchRequest) (*transport.SearchResponse, error) {
	log.Printf("Received: %v", in.GetPhonePart())
	p := cs.LoadPersistent()
	res, err := p.FindByPhone(in.GetPhonePart())
	if err != nil {
		return nil, errors.Wrap(err, "failed to find ")
	}
	response := &transport.SearchResponse{
		Result: make([]*transport.Person, 0, len(res)),
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
	transport.RegisterUserSearcherServer(s, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
