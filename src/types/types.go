package pgtypes

type Credentials struct {
	Username  string
	Password  string
	Database  string
	Host      string
	Port      int
	SSlConfig *SSL
}

type SSL struct {
	Certificate string
}

type operation int

// enum for pg wire protocol operations
const (
	QUERY operation = iota + 1
	STARTUP
)
