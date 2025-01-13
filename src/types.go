package src

type credentials struct {
	username string
	password string
	host     string
	port     int
}

type ssl_credentials struct {
	username string
	password string
	host     string
	port     int
	ssl      *ssl
}

type ssl struct {
	certificate string
	key         string
}
