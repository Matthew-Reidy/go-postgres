package types

type Credentials struct {
	Username           string
	Password           string
	Database           string
	Host               string
	Port               int
	InsecureSkipVerify bool
	SSlConfig          *SSL
}

type SSL struct {
	// optional path to a PEM bundle; when empty (or SSlConfig is nil)
	// the operating system's trusted cert chain is used
	Certificate string
}

type Operation int

// enum for pg wire protocol operations
const (
	QUERY Operation = iota + 1
	STARTUP
)
