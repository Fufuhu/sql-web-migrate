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

// TestBuildConnectionStringForUnixDomainSocket
// UnixDomainSocket用のPostgreSQL接続文字列が正常に生成できることを確認する。
func TestBuildConnectionStringForUnixDomainSocket(t *testing.T) {

	socketDirectoryPath := "/etc/postgres"
	user := "user"
	password := "password"
	dbname := "dbname"

	connectionString := BuildConnectionStringForUnixDomainSocket(
		socketDirectoryPath,
		user,
		password,
		dbname)

	if connectionString != "host=/etc/postgres user=user password=password dbname=dbname" {
		t.Log(connectionString)
		t.Fail()
	}
}
