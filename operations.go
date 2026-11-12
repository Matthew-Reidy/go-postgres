package postgres

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	msg "github.com/Matthew-Reidy/go-postgres/internal/message"
	pgtypes "github.com/Matthew-Reidy/go-postgres/types"
)

const (
	PG_SSL_SUPPORTED   byte = 0x53
	PG_SSL_UNSUPPORTED byte = 0x4E
)

// Conn is an established, TLS protected connection to a postgres server.
type Conn struct {
	*tls.Conn
}

// establishes as standard tcp connection to the server
// go-postgres then begins start up routine with the postgres server
func Connect(connectionConfig *pgtypes.Credentials) (*Conn, error) {

	if connectionConfig == nil {
		return nil, fmt.Errorf("connection config is nil")
	}

	host := fmt.Sprintf("%v:%d", connectionConfig.Host, connectionConfig.Port)

	dialer := &net.Dialer{
		Timeout: time.Second * 30,
	}
	conn, err := dialer.Dial("tcp", host)

	if err != nil {
		return nil, err
	}
	//check if the server supports a ssl/tls connection
	ssl_req := []byte{0x0, 0x0, 0x0, 0x8, 0x4, 0xd2, 0x16, 0x2f}

	if _, err := conn.Write(ssl_req); err != nil {
		conn.Close()
		return nil, fmt.Errorf("sending SSL request: %w", err)
	}

	ssl_resp := make([]byte, 1)

	if _, err := io.ReadFull(conn, ssl_resp); err != nil {
		conn.Close()
		return nil, fmt.Errorf("reading SSL response: %w", err)
	}

	//if ssl is not supported abort TLS connection attempt. Intended for localhost testing only. Not recommended for production use.
	if ssl_resp[0] != PG_SSL_SUPPORTED {
		log.Println("server does not support SSL/TLS connection...aborting TLS attempt")
		return nil, nil
	}

	log.Println("ssl connection allowed...begining handshake")

	tlsConn, err := establishTLSConnection(conn, connectionConfig)

	if err != nil {
		conn.Close()
		return nil, err
	}

	startup, err := msg.Encode(connectionConfig, "startup")

	if err != nil {
		tlsConn.Close()
		return nil, err
	}

	if _, err := tlsConn.Write(startup); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("sending startup message: %w", err)
	}

	resp := make([]byte, 42)

	if _, err := tlsConn.Read(resp); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("reading startup response: %w", err)
	}

	msg.Decode("AuthOptions", resp)

	return &Conn{tlsConn}, nil
}

// performs tls handshake
func establishTLSConnection(conn net.Conn, credentials *pgtypes.Credentials) (*tls.Conn, error) {

	tlsConn := tls.Client(conn, &tls.Config{ServerName: credentials.Host,
		RootCAs:    nil,
		MaxVersion: tls.VersionTLS12,
		MinVersion: tls.VersionTLS12})

	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("tls handshake: %w", err)
	}

	return tlsConn, nil
}

// sends query to the postgres server and recieves a response
func (conn *Conn) Query(queryString string) ([]byte, error) {

	if conn == nil || conn.Conn == nil {
		return nil, fmt.Errorf("connection is not open")
	}

	message, err := msg.Encode(queryString, pgtypes.QUERY)

	if err != nil {
		return nil, err
	}

	if _, err := conn.Write(message); err != nil {
		return nil, err
	}

	return []byte{}, nil
}

// Disconnect is safe to call on a nil or failed connection.
func (conn *Conn) Disconnect() {

	if conn == nil || conn.Conn == nil {
		return
	}

	//terminates the session database side
	b := []byte{'X', 0x00, 0x00, 0x00, 0x04}

	conn.Write(b)

	//terminates the actual connection
	if err := conn.Close(); err != nil {
		log.Println("error closing connection:", err)
	}
}
