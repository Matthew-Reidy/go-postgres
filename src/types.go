package operations

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
