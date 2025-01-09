package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func loadCertificates() (*x509.CertPool, error) {
	file, err := os.ReadFile("")

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

func connect(connConfig *connbject) (*tls.Conn, error) {

	certificates, err := loadCertificates()

	if err != nil {
		log.Fatal("FATAL ERROR!", err)
	}

	conn, err := tls.Dial("tcp", connConfig.host, &tls.Config{RootCAs: certificates})

	if err != nil {
		log.Fatal("FATAL ERROR!", err)
	}

	return conn, err
}
