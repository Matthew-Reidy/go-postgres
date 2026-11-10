package authentication

import (
	"fmt"

	"github.com/xdg-go/scram"
)

func SCRAM256(password string, username string, salt []byte, iterations int) (string, error) {

	client, err := scram.SHA256.NewClient(username, password, "")

	if err != nil {
		fmt.Println("Error creating SCRAM client:", err)
	}

	conv := client.NewConversation()

}

func sendClientMessage() {

}

func recieveServerMessage() {

}
