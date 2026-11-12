package authentication

import (
	"crypto/md5"
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

func generateSalt() string {
	return ""
}
