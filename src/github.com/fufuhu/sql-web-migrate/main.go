package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	config "github.com/fufuhu/sql-web-migrate/migrate"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

func main() {

	// URLパスと関数の関係を定義
	http.HandleFunc("/migrate/up", execMigrateUp)
	http.HandleFunc("/migrate/down", execMigrateDown)

	// ListenするIPアドレスを定義
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		fmt.Println(err)
	}

}

// LogRecord DDLの適用記録を格納するための構造体
// ID SQLのID
// AppliedAt SQLの適用タイムスタンプ
type LogRecord struct {
	ID        string `json:id`
	AppliedAt string `json:appliedAt`
}

const (
	// ErrorKey Logger output key attribure
	ErrorKey = "err"
)

// 一旦Postgre固定
func getConnection(connectionConfig config.DBConnectionConfig, dialect string) (*sql.DB, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	host := connectionConfig.Host()
	port, err := connectionConfig.Port()
	if err != nil {
		logger.Error(
			"Failed to get TCP port config",
			zap.Error(err))
	}
	sslMode, err := connectionConfig.SSLMode()
	if err != nil {
		logger.Error(
			"Failed to get SSL mode config",
			zap.Error(err))
	}

	var connectionString string
	if host[:1] == "/" {
		connectionString = config.BuildConnectionStringForUnixDomainSocket(
			host,
			connectionConfig.User(),
			connectionConfig.Password(),
			connectionConfig.DBName())
	} else {
		connectionString = config.BuildConnectionString(
			host,
			port,
			connectionConfig.User(),
			connectionConfig.Password(),
			connectionConfig.DBName(),
			sslMode)
	}

	db, err := sql.Open(config.DialectPostgres, connectionString)

	return db, err
}

func execMigrate(direction migrate.MigrationDirection) ([]LogRecord, error) {

	sourcePath := config.GetMigrationSourcePath()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info(
		"Setup source file path to migrate",
		zap.String("sourcePath", sourcePath))

	source := migrate.FileMigrationSource{
		Dir: sourcePath,
	}

	db, err := getConnection(
		config.ConnectionConfig,
		config.DialectPostgres)

	if err != nil {
		logger.Error(
			"DB connection open failure",
			zap.Error(err))
	}

	n, err := migrate.ExecMax(db, config.DialectPostgres, source, direction, 0)
	if err != nil {
		logger.Error(
			"Migration failed",
			zap.Error(err))
	}

	logger.Info(fmt.Sprintf("Applied %d migrations!", n))

	// クエリを投げる

	statement, err := db.Prepare("select id, applied_at from gorp_migrations")

	if err != nil {
		logger.Error(
			"Result check failed",
			zap.Error(err))
	}
	defer statement.Close()

	var (
		id        string
		appliedAt time.Time
	)

	rows, err := statement.Query()

	// 結果の取得と構造体への格納
	records := []LogRecord{}
	for rows.Next() {
		err := rows.Scan(&id, &appliedAt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, appliedAt)
		record := LogRecord{
			ID:        id,
			AppliedAt: appliedAt.String(),
		}

		records = append(records, record)
	}

	return records, err
}

const (
	// XForwardedFor X-Forwarded-For constant for IDE
	XForwardedFor = "X-Forwarded-For"
	// RemoteAddr RemoteAddr constant for IDE
	RemoteAddr = "RemoteAddr"
)

func getXForwardedFor(r *http.Request) string {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	addresses := r.Header.Get(XForwardedFor)
	logger.Info(
		"Checking HTTP Request Header, X-Forwarded-For",
		zap.String(XForwardedFor, addresses),
	)

	return addresses
}

func getRemoteAddr(r *http.Request) string {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	remote := r.RemoteAddr
	logger.Info(
		"Checking HTTP Request Header, RemoteAddr",
		zap.String("RemoteAddr", remote),
	)
	return remote
}

func execMigrateUp(w http.ResponseWriter, r *http.Request) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// addresses check

	networks := config.NetworkConfig.AllowedNetworks()

	addresses := strings.Split(getXForwardedFor(r), " ")
	remoteAddr := strings.Split(getRemoteAddr(r), ":")[0]
	addresses = append(addresses, remoteAddr)

	var ips []net.IP
	for _, address := range addresses {
		ips = append(ips, net.ParseIP(address))
	}

	isAllowed := false
	for _, ip := range ips {
		isAllowed = networks.IsAllowed(ip)
		if isAllowed {
			break
		}
	}

	if !isAllowed {
		logger.Warn(
			"Access forbidden",
		)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// migrationの実行
	records, err := execMigrate(migrate.Up)
	if err != nil {
		logger.Error(
			"Migration failed",
			zap.Error(err))
	}

	bytes, _ := json.Marshal(records)

	recordString := string(bytes)

	fmt.Fprintf(w, "%s\n", recordString)

}

func execMigrateDown(w http.ResponseWriter, r *http.Request) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// addresses check

	networks := config.NetworkConfig.AllowedNetworks()

	addresses := strings.Split(getXForwardedFor(r), " ")
	remoteAddr := strings.Split(getRemoteAddr(r), ":")[0]
	addresses = append(addresses, remoteAddr)

	var ips []net.IP
	for _, address := range addresses {
		ips = append(ips, net.ParseIP(address))
	}

	isAllowed := false
	for _, ip := range ips {
		isAllowed = networks.IsAllowed(ip)
		if isAllowed {
			break
		}
	}

	if !isAllowed {
		logger.Warn(
			"Access forbidden",
		)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// migrationの実行

	records, err := execMigrate(migrate.Down)
	if err != nil {
		logger.Error(
			"Migration failed",
			zap.Error(err))
	}

	bytes, _ := json.Marshal(records)

	recordString := string(bytes)

	fmt.Fprintf(w, "%s\n", recordString)

}
