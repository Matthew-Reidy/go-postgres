package src

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func loadCertificates(certPath string) (*x509.CertPool, error) {
	file, err := os.ReadFile(certPath)

	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()

	ok := roots.AppendCertsFromPEM(file)

	if !ok {

		return nil, err

	}

	return roots, nil

}

// establishes an ssl connection to the postgres server
// go-postgres then begins start up routine with the postgres server
func ssl_connect(connConfig *ssl_credentials) {

	host := fmt.Sprintf("%v:%d", connConfig.host, connConfig.port)

	certificates, err := loadCertificates(connConfig.ssl.certificate)

	if err != nil {
		log.Fatal("FATAL!", err)
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{RootCAs: certificates})

	defer conn.Close()

	if err != nil {
		log.Fatal("FATAL!", err)
	}

	defer conn.Close()

}

// establishes as standard tcp connection to the server
// go-postgres then begins start up routine with the postgres server
func connect(connConfig *credentials) {

	host := fmt.Sprintf("%v:%d", connConfig.host, connConfig.port)

	conn, err := net.Dial("tcp", host)

	if err != nil {
		log.Fatalf("FATAL! %v", err)
	}

	defer conn.Close()

}

func query(queryString string) {

}
