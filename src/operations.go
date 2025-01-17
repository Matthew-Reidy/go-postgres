package operations

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Matthew-Reidy/go-postgres/src/utils"
)

const (
	PG_SSL_SUPPORTED   = 0x53
	PG_SSL_UNSUPPORTED = 0x4E
)

// enum for pg wire protocol operations
const (
	Query = iota + 1
)

type tlsConnection struct {
	*tls.Conn
}

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

// establishes as standard tcp connection to the server
// go-postgres then begins start up routine with the postgres server
func Connect(connectionConfig *Credentials) (*tlsConnection, error) {

	host := fmt.Sprintf("%v:%d", connectionConfig.Host, connectionConfig.Port)

	dialer := &net.Dialer{
		Timeout: time.Second * 30,
	}
	conn, err := dialer.Dial("tcp", host)

	if err != nil {
		log.Fatal(err)
		return nil, err

	}

	ssl_req := []byte{0x0, 0x0, 0x0, 0x8, 0x4, 0xd2, 0x16, 0x2f}

	conn.Write(ssl_req)

	ssl_resp := make([]byte, 1)

	conn.Read(ssl_resp)

	log.Println(ssl_resp)

	//if ssl is supported by the server begin tls handshake
	if ssl_resp[0] == PG_SSL_SUPPORTED {

		log.Println("ssl connection allowed...begining handshake")

		tlsConn := establishTLSConnection(&conn, connectionConfig)

		//placeholder for now
		start_up := []byte{0x00, 0x00, 0x00, 0x17, 0x00, 0x03, 0x00, 0x00, 'u', 's', 'e', 'r', 0x00, 'p', 'o', 's', 't', 'g', 'r', 'e', 's', 0x00, 0x00}

		b := make([]byte, 150)

		tlsConn.Write(start_up)

		tlsConn.Read(b)

		log.Println(b)

		return &tlsConnection{tlsConn}, nil

	}

	//if ssl is not supported close the connection
	conn.Close()
	//and return error...
	return nil, fmt.Errorf("server does not support SSL/TLS connection...Aborting")
}

// performs tls handshake
func establishTLSConnection(conn *net.Conn, credentials *Credentials) *tls.Conn {

	cert, err := loadCertificates(credentials.SSlConfig.Certificate)

	if err != nil {
		log.Fatalf("FATAL! CANT LOAD CERTIFICATE: %x\n", err)
	}

	tlsConn := tls.Client(*conn, &tls.Config{ServerName: credentials.Host,
		RootCAs:    cert,
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13})

	return tlsConn
}

// sends query to the postgres server and recieves a response
func (conn *tlsConnection) Query(queryString string) ([]byte, error) {
	operation := "Q"
	message := utils.Encode(operation + queryString)

	_, err := conn.Write(message)

	if err != nil {
		log.Panic(err)
		return nil, err

	}

	return []byte{}, nil
}

func (conn *tlsConnection) Disconnect() {

	//terminates the session server side
	b := []byte{'X', 0x00, 0x00, 0x00, 0x05}

	conn.Write(b)

	//terminates the actual connection
	err := conn.Close()

	if err != nil {
		log.Panic(err)
	}

}
