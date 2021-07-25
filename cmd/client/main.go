package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/wshaman/course-grpc/common/transport"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	address       = "localhost:50051"
	defaultSearch = "call-me"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := transport.NewUserSearcherClient(conn)

	// Contact the server and print out its response.
	q := defaultSearch
	if len(os.Args) > 1 {
		q = os.Args[1]
	}

	for {
		caller(c, q)
	}

}

func caller(c transport.UserSearcherClient, q string) {
	fmt.Println("Input phone part")
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	pp, _ := reader.ReadString('\n')
	pp = strings.Replace(pp, "\n", "", -1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Search(ctx, &transport.SearchRequest{
		PhonePart: pp,
	})
	if err != nil {
		log.Fatalf("could not find: %v", err)
	}
	result := r.GetResult()
	log.Printf("Found: %n results\n", len(result))
	for _, v := range result {
		log.Println(v.String())
	}
}
