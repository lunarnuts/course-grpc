package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wshaman/course-grpc/common/transport"

	"google.golang.org/grpc"
)

const (
	address             = "localhost:50051"
	defaultUserRegister = "Constantine Phelps"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := transport.NewUserRegisterClient(conn)

	// Contact the server and print out its response.
out:
	for {
		fmt.Println("Choose action:")
		fmt.Println("1 -> Register")
		fmt.Println("2 -> List all users")
		fmt.Println("3 -> Close")
		fmt.Print("->")
		reader := bufio.NewReader(os.Stdin)
		pp, _ := reader.ReadString('\n')
		pp = strings.Replace(pp, "\n", "", -1)
		switch pp {
		case "1":
			Register(c)
		case "2":
			List(c)
		case "3":
			break out
		default:
			log.Print("Enter valid input")
		}
	}
}

func Register(c transport.UserRegisterClient) {
	fmt.Println("Input user name to register")
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	pp, _ := reader.ReadString('\n')
	pp = strings.Replace(pp, "\n", "", -1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UserRegister(ctx, &transport.UserRegisterRequest{
		Name: pp,
	})
	if err != nil {
		log.Printf("could not register: %v", err)
		return
	}
	result := r.GetResult()
	log.Printf("User registered with id: %d\n", result)
}

func List(c transport.UserRegisterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Println("All users list:")
	list, err := c.UsersList(ctx, &transport.EmptyRequest{})
	if err != nil {
		log.Printf("could not list: %v", err)
	}
	for _, v := range list.GetResult() {
		log.Printf("id: %d  name: \"%s\"\n", v.Id, v.Name)
	}
}
