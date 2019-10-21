package migrate

import (
	"errors"
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
	// DBMigrationSourcePath DBのマイグレーションファイルの含まれているディレクトリパス
	DBMigrationSourcePath = "SQL_MIGRATE_MIGRATION_SOURCE_PATH"
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
	DefaultDBSSLMode = "disable"
	// DefaultDBMigrationSourcePath デフォルトのマイグレーションのソースパス
	DefaultDBMigrationSourcePath = "/etc/migrate"
)

// GetHost DBホスト名を取得する。
// 環境変数が設定されていない場合は、DefaultDBHostの値を返す
func GetHost() string {
	return getValue(DBHost, DefaultDBHost)
}

// GetPort DBのTCPポート番号を取得する。
// 環境変数が設定されていない場合は、DefaultDBPortの値を返す
// 環境変数に数字以外が設定されている場合はerrorとポート番号として-1を返す
func GetPort() (int, error) {
	portString := getValue(DBPort, strconv.Itoa(DefaultDBPort))
	port, err := strconv.Atoi(portString)

	if err != nil {
		fmt.Println(err) //あとでzapに変えましょうね。
		return -1, err
	}
	return port, err
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

const (
	//SSLModeSettingFormatErrorMessage SSLModeの設定を誤っている際のエラーメッセージです
	SSLModeSettingFormatErrorMessage = "SSLMode should be require, verify-full, verify-ca, or disable"
)

// GetSSLMode SSLModeの有効/無効を取得する。
// 環境変数が設定されていない場合は、DefaultDBSSLModeの値を返す
// 不正な値(true/false以外)が設定されている場合はエラーとDefaultDBSSLModeの値を返す
func GetSSLMode() (string, error) {
	mode := getValue(DBSSLMode, DefaultDBSSLMode)

	if mode != "require" && mode != "verify-full" && mode != "verify-ca" && mode != "disable" {
		return DefaultDBSSLMode, errors.New(SSLModeSettingFormatErrorMessage)
	}
	return mode, nil
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

// GetMigrationSourcePath DBマイグレーション用のSQLファイルを格納しているディレクトリパス
func GetMigrationSourcePath() string {
	return getValue(DBMigrationSourcePath, DefaultDBMigrationSourcePath)
}

// DBConnectionConfig DBの接続情報を格納したもの
type DBConnectionConfig struct {
	Host     func() string
	Port     func() (int, error)
	User     func() string
	Password func() string
	DBName   func() string
	SSLMode  func() (string, error)
}

// ConnectionConfig Databaseへの接続設定です
var ConnectionConfig DBConnectionConfig

func init() {
	ConnectionConfig = DBConnectionConfig{
		Host:     GetHost,
		Port:     GetPort,
		User:     GetUser,
		Password: GetPassword,
		DBName:   GetDBName,
		SSLMode:  GetSSLMode,
	}
}
