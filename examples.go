package main

import (
	"log"

	operations "github.com/Matthew-Reidy/go-postgres/src"
)

func main() {

	connConfig := &operations.Credentials{
		Username: "hello",
		Password: "world",
		Host:     "myhost",
		Port:     5432,
		SSlConfig: &operations.SSL{
			Certificate: "some/path/way/cert.pem",
		},
	}

	conn, err := operations.Connect(connConfig)

	defer conn.Close()

	if err != nil {
		log.Fatalf("FATAL! : %x", err)
	}

	data, err := conn.Query("select * from users;")

	if err != nil {
		log.Panic(err)
	}

	//...do something with the data WIP

}
