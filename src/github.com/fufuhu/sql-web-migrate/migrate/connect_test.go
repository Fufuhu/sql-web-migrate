package migrate

import (
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	host := "host"
	port := 1234
	user := "user"
	password := "password"
	dbname := "dbname"
	sslmode := "false"

	connectionString := BuildConnectionString(host, port, user, password, dbname, sslmode)

	if connectionString != "host=host port=1234 user=user password=password dbname=dbname sslmode=false" {
		t.Log(connectionString)
		t.Fail()
	}
}
