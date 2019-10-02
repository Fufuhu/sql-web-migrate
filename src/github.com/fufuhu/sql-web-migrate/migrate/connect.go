package migrate

import "fmt"

const (
	// ConnectionStringTemplate 接続文字列のテンプレート
	ConnectionStringTemplate = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
)

// BuildConnectionString SQLの接続文字列を生成する
func BuildConnectionString(host string, port int, user string,
	password string, dbname string, sslmode string) string {

	connectionString := fmt.Sprintf(ConnectionStringTemplate, host, port, user, password, dbname, sslmode)

	return connectionString
}
