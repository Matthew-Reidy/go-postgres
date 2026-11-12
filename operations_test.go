package postgres_test

import (
	"testing"

	pg "github.com/Matthew-Reidy/go-postgres"
	pgtypes "github.com/Matthew-Reidy/go-postgres/types"
)

func TestConnection(t *testing.T) {

	// connConfig := &pgtypes.Credentials{
	// 	Username: "postgres",
	// 	Password: "mypw",
	// 	Database: "postgres",
	// 	Host:     "mytestdb.ceg3kwpf6czu.us-west-1.rds.amazonaws.com",
	// 	Port:     5432,
	// 	SSlConfig: &pgtypes.SSL{
	// 		Certificate: "global-bundle.pem",
	// 	},
	// }

	connConfig := &pgtypes.Credentials{
		Username: "postgres",
		Password: "mypw",
		Database: "postgres",
		Host:     "localhost",
		Port:     5432,
	}

	_, err := pg.Connect(connConfig)

	if err != nil {
		t.Fatalf("An error occured! %s", err)
	}

}
