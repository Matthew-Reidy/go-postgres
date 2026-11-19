package postgres_test

import (
	"testing"

	pg "github.com/Matthew-Reidy/go-postgres/internal"
	"github.com/Matthew-Reidy/go-postgres/types"
)

func TestConnectionUnsafe(t *testing.T) {

	connConfig := &types.Credentials{
		Username: "postgres",
		Password: "postgres1234",
		Database: "postgres",
		Host:     "localhost",
		Port:     5432,
	}

	_, err := pg.Connect(connConfig)

	if err != nil {
		t.Fatalf("An error occured! %s", err)
	}

}

func TestConnectionSafe(t *testing.T) {

	connConfig := &types.Credentials{
		Username: "postgres",
		Password: "postgres1234",
		Database: "postgres",
		Host:     "pg-test.ceg3kwpf6czu.us-west-1.rds.amazonaws.com",
		Port:     5432,
	}

	_, err := pg.Connect(connConfig)

	if err != nil {
		t.Fatalf("An error occured! %s", err)
	}

}

func TestConnectionWithInvalidCredentials(t *testing.T) {
	connConfig := &types.Credentials{
		Username: "invalid_user",
		Password: "invalid_password",
		Database: "postgres",
		Host:     "localhost",
		Port:     5432,
	}

	_, err := pg.Connect(connConfig)

	if err == nil {
		t.Fatalf("Expected an error due to invalid credentials, but got none.")
	}
}

func testQuery(t *testing.T) {

}
