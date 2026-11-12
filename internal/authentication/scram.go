package authentication

import (
	"fmt"

	"github.com/xdg-go/scram"
)

func SCRAM256(password string, username string, salt []byte, iterations int) (string, error) {

	client, err := scram.SHA256.NewClient(username, password, "")

	if err != nil {
		return "", fmt.Errorf("error creating SCRAM client: %w", err)
	}

	conv := client.NewConversation()
	_ = conv // TODO: drive the SCRAM exchange with conv

	return "", fmt.Errorf("SCRAM-SHA-256 not implemented yet")

}

func sendClientMessage() {

}

func recieveServerMessage() {

}
