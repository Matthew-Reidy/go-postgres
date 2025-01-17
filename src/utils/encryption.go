package utils

import (
	"crypto/md5"
)

func MD5(password string, salt any) [md5.Size]byte {

	return md5.Sum([]byte(password))

}
