package message

import (
	"fmt"

	"github.com/Matthew-Reidy/go-postgres/types"
)

const (
	AUTH_RESP = 'p'
)

// Encode builds the wire-format bytes for a message. It performs no I/O,
// so callers own the connection and are free to use it concurrently.
func Encode(message any, messageType types.Operation) ([]byte, error) {
	messageBody := []byte{}
	switch messageType {

	case types.STARTUP:

		creds, ok := message.(*types.Credentials)

		if !ok || creds == nil {
			return nil, fmt.Errorf("startup message requires *types.Credentials")
		}

		return nil, nil

	case types.QUERY:

		query, ok := message.(string)

		if !ok || len(query) == 0 {
			return nil, fmt.Errorf("query message requires a string")
		}

		return nil, nil

	}

	return nil, fmt.Errorf("unknown message type %q", messageType)
}

func MD5PasswordRoutine() {

}

func SCRAMRoutine() {
	fmt.Println("starting scram")
}
