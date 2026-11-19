package postgres

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	msg "github.com/Matthew-Reidy/go-postgres/internal/message"
	"github.com/Matthew-Reidy/go-postgres/types"
)

const (
	PG_SSL_SUPPORTED   byte = 0x53
	PG_SSL_UNSUPPORTED byte = 0x4E
)

// Pgconn is an established connection to a postgres server, either the raw
// TCP connection or a TLS connection wrapping it.
type Pgconn struct {
	net.Conn
}

// establishes as standard tcp connection to the server
// go-postgres then begins start up routine with the postgres server
func Connect(connectionConfig *types.Credentials) (*Pgconn, error) {

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
	fmt.Printf("ssl_req %x", ssl_req)

	if _, err := conn.Write(ssl_req); err != nil {
		conn.Close()
		return nil, fmt.Errorf("sending SSL request: %w", err)
	}

	ssl_resp := make([]byte, 1)

	if _, err := io.ReadFull(conn, ssl_resp); err != nil {

		conn.Close()
		return nil, fmt.Errorf("reading SSL response: %w", err)
	}
	fmt.Printf("ssl_resp %s", string(ssl_resp))
	//if ssl is not supported abort TLS connection attempt. Intended for localhost testing only. Not recommended for production use.
	if ssl_resp[0] != PG_SSL_SUPPORTED || connectionConfig.InsecureSkipVerify == true {
		log.Println("server does not support SSL/TLS connection...aborting TLS attempt")
		return &Pgconn{conn}, nil
	}

	log.Println("ssl connection allowed...begining handshake")

	tlsConn, err := establishTLSConnection(conn, connectionConfig)

	if err != nil {
		conn.Close()
		return nil, err
	}

	startup, err := msg.Encode(connectionConfig, types.STARTUP)

	if err != nil {
		tlsConn.Close()
		return nil, err
	}

	if _, err := tlsConn.Conn.Write(startup); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("sending startup message: %w", err)
	}

	resp := make([]byte, 42)

	if _, err := tlsConn.Conn.Read(resp); err != nil {
		fmt.Println(resp)
		tlsConn.Close()
		return nil, fmt.Errorf("reading startup response: %w", err)
	}

	msg.Decode("AuthOptions", resp)

	return tlsConn, nil
}

// performs tls handshake
func establishTLSConnection(conn net.Conn, credentials *types.Credentials) (*Pgconn, error) {

	tlsConn := tls.Client(conn, &tls.Config{ServerName: credentials.Host,
		RootCAs:            nil, //default is nil so we're using nil
		MaxVersion:         tls.VersionTLS12,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true})

	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("tls handshake: %w", err)
	}

	return &Pgconn{tlsConn}, nil
}

// accepts raw connection or tls connection
func (conn *Pgconn) Startup() ([]byte, error) {
	return nil, nil
}

// sends query to the postgres server and recieves a response
func (conn *Pgconn) Query(queryString string) ([]byte, error) {

	if conn == nil {
		return nil, fmt.Errorf("connection is not open")
	}

	message, err := msg.Encode(queryString, types.QUERY)

	if err != nil {
		return nil, err
	}

	if _, err := conn.Write(message); err != nil {
		return nil, err
	}

	return []byte{}, nil
}

// Disconnect is safe to call on a nil or failed connection.
func (conn *Pgconn) Disconnect() {

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
