package main

import (
	"log"

	pg "github.com/Matthew-Reidy/go-postgres/src"
)

func main() {

	connConfig := &pg.Credentials{
		Username: "hello",
		Password: "world",
		Database: "mydb",
		Host:     "myhost",
		Port:     5432,
		SSlConfig: &pg.SSL{
			Certificate: "some/path/way/cert.pem",
		},
	}

	conn, err := pg.Connect(connConfig)

	defer conn.Close()

	if err != nil {
		log.Fatalf("FATAL! : %x", err)
	}

	data, err := conn.Query("select * from users;")

	if err != nil {
		log.Panic(err)
	}

	//...do something with the data WIP

	conn.Disconnect()

}
