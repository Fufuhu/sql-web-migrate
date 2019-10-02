package migrate

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// DBHost DBホストを指定するための環境変数
	DBHost = "SQL_MIGRATE_HOST"
	// DBPort DBがListenしているTCPポート番号の環境変数
	DBPort = "SQL_MIGRATE_PORT"
	// DBUser DBのユーザを指定するための環境変数
	DBUser = "SQL_MIGRATE_USER"
	// DBPassword DBのパスワードを指定するための環境変数
	DBPassword = "SQL_MIGRATE_PASSWORD"
	// DBName DBの名前を指定するための環境変数
	DBName = "SQL_MIGRATE_DBNAME"
	// DBSSLMode SSLモードの有効・向こうを指定するための環境変数
	DBSSLMode = "SQL_MIGRATE_SSL_MODE"
)

const (
	// DefaultDBHost デフォルトのDBホスト
	DefaultDBHost = "localhost"
	// DefaultDBPort DBがListenしているデフォルトのTCPポート番号
	DefaultDBPort = 5432
	// DefaultDBUser デフォルトのDBユーザ
	DefaultDBUser = ""
	// DefaultDBPassword デフォルトのDBパスワード
	DefaultDBPassword = ""
	// DefaultDBName デフォルトのDB名
	DefaultDBName = ""
	// DefaultDBSSLMode デフォルトのSSLモード設定
	DefaultDBSSLMode = "false"
)

// GetHost DBホスト名を取得する。
// 環境変数が設定されていない場合は、DefaultDBHostの値を返す
func GetHost() string {
	return getValue(DBHost, DefaultDBHost)
}

// GetPort DBのTCPポート番号を取得する。
// 環境変数が設定されていない場合は、DefaultDBPortの値を返す
func GetPort() int {
	portString := getValue(DBPort, strconv.Itoa(DefaultDBPort))
	port, err := strconv.Atoi(portString)

	if err != nil {
		fmt.Println(err) //あとでzapに変えましょうね。
	}
	return port
}

// GetUser DBのユーザ名を取得する
// 環境変数が設定されていない場合は、DefaultDBUserの値を返す
func GetUser() string {
	return getValue(DBUser, DefaultDBUser)
}

// GetPassword DBのパスワードを取得する
// 環境変数が設定されていない場合は、DefaultPasswordの値を返す
func GetPassword() string {
	return getValue(DBPassword, DefaultDBPassword)
}

// GetDBName DBの名前を取得する
// 環境変数が設定されていない場合は、DefaultDBNameの値を返す
func GetDBName() string {
	return getValue(DBName, DefaultDBName)
}

// GetSSLMode SSLModeの有効/無効を取得する。
// 環境変数が設定されていない場合は、DefaultDBSSLModeの値を返す
func GetSSLMode() string {
	return getValue(DBSSLMode, DefaultDBSSLMode)
}

// getValue envKeyに指定された環境変数の値を返す。
// envKeyに指定された環境変数が未定義だった場合はdefaultValueを返す。
func getValue(envKey string, defaultValue string) string {
	value := os.Getenv(envKey)

	if value == "" {
		value = defaultValue
	}

	return value
}