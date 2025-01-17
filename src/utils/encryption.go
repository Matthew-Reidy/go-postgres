package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

/*
hashes the password with the md5 hashing algorithm as described inthe Postgres docs

https://www.postgresql.org/docs/current/protocol-flow.html#PROTOCOL-FLOW-START-UP

concat('md5', md5(concat(md5(concat(password, username)), random-salt)))
*/
func MD5HashPassword(password string, username string, salt string) string {

	firstConcat := fmt.Sprintf("%s%s", password, username)

	firstHash := md5.Sum([]byte(firstConcat))

	firstHashHex := hex.EncodeToString(firstHash[:])

	secondConcat := firstHashHex + salt
	secondHash := md5.Sum([]byte(secondConcat))

	return fmt.Sprintf("md5%s", hex.EncodeToString(secondHash[:]))

}

// for use in SCRAM-SHA-256 authentication in the client-first step
func generateClientNonce() (string, error) {

	bytes := make([]byte, 10)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

// client proof as described by
func SCRAMKeys(password string, salt string) (string, string) {
	clientKey := ""
	serverKey := ""

	key := []byte(password + salt)

	hmac := hmac.New(sha256.New, key)

	for i := range 2 {
		if i%2 == 0 {
			hmac.Write([]byte("Client Key"))

			clientKey = string(hmac.Sum(nil))

			hmac.Reset()

			continue
		}
		hmac.Write([]byte("Server Key"))

		clientKey = string(hmac.Sum(nil))

		hmac.Reset()

	}

	return clientKey, serverKey
}

func generateAuthValue(clientServerFirsts *ServerClientMessage) []byte {

}
