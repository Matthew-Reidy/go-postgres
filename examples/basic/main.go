package main

import (
	"log"

	pg "github.com/Matthew-Reidy/go-postgres/internal"
	"github.com/Matthew-Reidy/go-postgres/types"
)

func main() {

	connConfig := &types.Credentials{
		Username: "hello",
		Password: "world",
		Database: "mydb",
		Host:     "myhost",
		Port:     5432,
		SSlConfig: &types.SSL{
			Certificate: "some/path/way/cert.pem",
		},
	}

	conn, err := pg.Connect(connConfig)

	if err != nil {
		log.Fatalf("FATAL! : %v", err)
	}

	defer conn.Disconnect()

	_, err = conn.Query("select * from my_table;")

	if err != nil {
		log.Panic(err)
	}

	//...do something with the data (WIP)

}
