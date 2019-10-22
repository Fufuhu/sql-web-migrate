package migrate

import "fmt"

const (
	// ConnectionStringTemplate 接続文字列のテンプレート
	ConnectionStringTemplate = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	// ConnectionStringForUnixDomainSocketTemplate Template string to use Unix domain socket
	ConnectionStringForUnixDomainSocketTemplate = "host=%s user=%s password=%s dbname=%s"
)

const (
	// DialectPostgres PostgreSQLを使う際の指定子
	DialectPostgres = "postgres"
)

// BuildConnectionString PostgreSQLの接続文字列を生成する
func BuildConnectionString(host string, port int, user string,
	password string, dbname string, sslmode string) string {

	connectionString := fmt.Sprintf(ConnectionStringTemplate, host, port, user, password, dbname, sslmode)

	return connectionString
}

// BuildConnectionStringForUnixDomainSocket PostgreSQLの接続文字列を生成する
func BuildConnectionStringForUnixDomainSocket(socketDirectoryPath string, user string,
	password string, dbname string) string {

	connectionString := fmt.Sprintf(ConnectionStringForUnixDomainSocketTemplate,
		socketDirectoryPath, user, password, dbname)

	return connectionString
}
