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

func saltShaker(serverFirst *ServerFirstMessage, password string) string {
	saltedpw := []byte{}

	initHash := hmac.New(sha256.New, []byte(password))

	initHash.Write([]byte(serverFirst.Salt))

	saltedpw = initHash.Sum(nil)

	for i := 2; i <= serverFirst.IterationCount; i++ {
		iterHash := hmac.New(sha256.New, []byte(password))
		iterHash.Write(saltedpw)
		saltedpw = iterHash.Sum(nil)
	}

	return hex.EncodeToString(saltedpw)

}

// client proof as described by
func SCRAMClientKey(saltedPassword string) string {
	clientKey := []byte{}

	key := []byte(saltedPassword)

	hmac := hmac.New(sha256.New, key)

	hmac.Write([]byte("Client Key"))

	clientKey = hmac.Sum(nil)

	hmac.Reset()

	return hex.EncodeToString(clientKey)
}
