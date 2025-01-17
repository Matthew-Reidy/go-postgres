package operations_test

import (
	"testing"

	pg "github.com/Matthew-Reidy/go-postgres/src"
)

func TestConnection(t *testing.T) {

	connConfig := &pg.Credentials{
		Username: "postgres",
		Password: "mypw",
		Database: "mydb",
		Host:     "myhost",
		Port:     5432,
		SSlConfig: &pg.SSL{
			Certificate: "global-bundle.pem",
		},
	}

	_, err := pg.Connect(connConfig)

	if err != nil {
		t.Fatalf("An error occured! %x", err)
	}

}
