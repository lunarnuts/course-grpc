package db

import (
	"fmt"
	"log"

	"github.com/wshaman/course-grpc/common/transport"
)

var Data = Database{
	Map:  make(map[string]*transport.Person),
	List: make([]*transport.Person, 0),
	Len:  int64(0),
}

type Database struct {
	Map  map[string]*transport.Person
	List []*transport.Person
	Len  int64
}

func LoadPersistent() []*transport.Person {
	return Data.List
}

func CheckIfExists(person *transport.UserRegisterRequest) bool {
	_, ok := Data.Map[person.Name]
	return ok
}

func AddToDB(person *transport.UserRegisterRequest) error {
	if CheckIfExists(person) {
		return fmt.Errorf("user already exists")
	}
	newPerson := &transport.Person{
		Id:   Data.Len,
		Name: person.Name,
	}
	log.Printf("New entry: id: %d  name: \"%s\"", newPerson.Id, newPerson.Name)
	Data.List = append(Data.List, newPerson)
	Data.Map[person.GetName()] = newPerson
	Data.Len++
	return nil
}

func Populate() {
	var List = []*transport.UserRegisterRequest{
		{Name: "Egor Uliyanov"},
		{Name: "ilya Komarskih"},
		{Name: "Ivan Ivanov"},
		{Name: "Alex Fergusson"},
		{Name: "Julian Casablancas"},
		{Name: "Pavel Lyadov"},
	}
	for _, person := range List {
		AddToDB(person)
	}
}
