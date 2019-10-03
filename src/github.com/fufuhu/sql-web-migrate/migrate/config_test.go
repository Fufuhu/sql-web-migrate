package migrate

import (
	"os"
	"strconv"
	"testing"
)

// TestGetHost DBホストとして環境変数に
// 指定した値を取得できていることを確認する。
func TestGetHost(t *testing.T) {
	exptected := "hogehoge"
	os.Setenv(DBHost, exptected)

	host := GetHost()

	if host != exptected {
		t.Fail()
	}
}

// TestGetHostDefaultValue DBホストとして
// 環境変数が指定されていない場合に
// デフォルト値が取得できることを確認する。
func TestGetHostDefaultValue(t *testing.T) {
	expected := DefaultDBHost
	os.Unsetenv(DBHost)

	host := GetHost()

	if host != expected {
		t.Fail()
	}
}

// TestGetPort DBのポートとして環境変数に
// 指定した値を取得できていることを確認する。
func TestGetPort(t *testing.T) {
	expected := 10000
	os.Setenv(DBPort, strconv.Itoa(expected))

	port, err := GetPort()

	if err != nil {
		t.Fail()
	}

	if port != expected {
		t.Fail()
	}
}

// TestGetPortDefaultValue DBのポートとして
// 環境変数が指定されていない場合に
// デフォルト値が取得できることを確認する。
func TestGetPortDefaultValue(t *testing.T) {
	expected := DefaultDBPort
	os.Unsetenv(DBPort)

	port, err := GetPort()

	if err != nil {
		t.Fail()
	}

	if port != expected {
		t.Fail()
	}
}

func TestGetPortWrongTypeFormatError(t *testing.T) {
	os.Setenv(DBPort, "hogehoge")

	port, err := GetPort()

	if port != -1 {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

// TestGetUser DBのユーザとして環境変数に
// 指定した値を取得できていることを確認する。
func TestGetUser(t *testing.T) {
	expected := "fugafuga"
	os.Setenv(DBUser, expected)

	user := GetUser()

	if user != expected {
		t.Fail()
	}
}

// TestGetUserDefaultValue DBのユーザとして
// 環境変数が指定されていない場合に
// デフォルト値が取得できることを確認する。
func TestGetUserDefaultValue(t *testing.T) {
	expected := DefaultDBUser
	os.Unsetenv(DBUser)

	user := GetUser()

	if user != expected {
		t.Fail()
	}
}

// TestGetPassword DBのパスワードとして環境変数に
// 指定した値を取得できていることを確認する。
func TestGetPassword(t *testing.T) {
	expected := "fugafuga"
	os.Setenv(DBPassword, expected)

	password := GetPassword()

	if password != expected {
		t.Fail()
	}
}

// TestGetPasswordDefaultValue DBのパスワードとして
// 環境変数が指定されていない場合に
// デフォルト値が取得できることを確認する。
func TestGetPasswordDefaultValue(t *testing.T) {
	expected := DefaultDBPassword
	os.Unsetenv(DBPassword)

	password := GetPassword()

	if password != expected {
		t.Fail()
	}
}

// TestGetDBName DBの名前として環境変数に
// 指定した値を取得できていることを確認する。
func TestGetDBName(t *testing.T) {
	expected := "fugafuga"
	os.Setenv(DBName, expected)

	dbName := GetDBName()

	if dbName != expected {
		t.Fail()
	}
}

// TestGetDBNameDefaultValue DBのパスワードとして
// 環境変数が指定されていない場合に
// デフォルト値が取得できることを確認する。
func TestGetDBNameDefaultValue(t *testing.T) {
	expected := DefaultDBPassword
	os.Unsetenv(DBName)

	dbName := GetPassword()

	if dbName != expected {
		t.Fail()
	}
}

// TestGetSSLMode DBのSSLModeアクセスが有効かを確認する。
// trueの値がかえってくることを期待する。
func TestGetSSLMode(t *testing.T) {
	expected := "true"
	os.Setenv(DBSSLMode, expected)

	sslMode, err := GetSSLMode()

	if err != nil {
		t.Fail()
	}

	if sslMode != expected {
		t.Fail()
	}
}

// TestGetSSLModeDefaultValue DBのSSLModeとして
// デフォルトの値がかえってくることを確認する。
func TestGetSSLModeDefaultValue(t *testing.T) {
	os.Unsetenv(DBSSLMode)

	sslMode, err := GetSSLMode()

	if err != nil {
		t.Fail()
	}

	if sslMode != DefaultDBSSLMode {
		t.Fail()
	}
}

// TestGetSSLModeError DBのSSLModeとして
// 不正な値が設定されてる場合にエラーが出力されること、
// デフォルトの値がかえってくることを確認する。
func TestGetSSLModeError(t *testing.T) {
	os.Setenv(DBSSLMode, "hogehoge")

	sslMode, err := GetSSLMode()

	if err.Error() != SSLModeSettingFormatErrorMessage {
		t.Fail()
	}

	if sslMode != DefaultDBSSLMode {
		t.Fail()
	}
}

// TestGetValue getValue関数に指定した環境変数のキーと
// それに紐づく値がかえってくることを確認する
func TestGetValue(t *testing.T) {
	envKey := "hoge"
	envValue := "piyo"
	defaultValue := "fuga"
	os.Setenv(envKey, envValue)

	value := getValue(envKey, defaultValue)

	if envValue != value {
		t.Fail()
	}
}

// TestGetValueDefault getValue関数にて指定したキーに
// 対応する環境変数の値が設定されていない場合に、
// デフォルトの値がかえってくることを確認する。
func TestGetValueDefault(t *testing.T) {
	envKey := "hoge"
	defaultValue := "fuga"
	os.Unsetenv(envKey)

	value := getValue(envKey, defaultValue)

	if defaultValue != value {
		t.Fail()
	}
}
