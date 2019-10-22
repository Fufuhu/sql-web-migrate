package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "github.com/fufuhu/sql-web-migrate/migrate"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {

	// URLパスと関数の関係を定義
	http.HandleFunc("/migrate/up", execMigrateUp)
	http.HandleFunc("/migrate/down", execMigrateDown)

	// ListenするIPアドレスを定義
	err := http.ListenAndServe("0.0.0.0:18080", nil)

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

// 一旦Postgre固定
func getConnection(connectionConfig config.DBConnectionConfig, dialect string) (*sql.DB, error) {

	host := connectionConfig.Host()
	port, err := connectionConfig.Port()
	if err != nil {
		fmt.Println(err)
	}
	sslMode, err := connectionConfig.SSLMode()
	if err != nil {
		fmt.Println(err)
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

	db, err := sql.Open("postgres", connectionString)

	return db, err
}

func execMigrate(direction migrate.MigrationDirection) ([]LogRecord, error) {

	sourcePath := config.GetMigrationSourcePath()
	fmt.Println(sourcePath)

	source := migrate.FileMigrationSource{
		Dir: sourcePath,
	}

	db, err := getConnection(
		config.ConnectionConfig,
		"postgres")

	if err != nil {
		fmt.Println("DB connection open failure.")
		fmt.Println(err)
	}

	n, err := migrate.ExecMax(db, "postgres", source, direction, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)

	// クエリを投げる

	statement, err := db.Prepare("select id, applied_at from gorp_migrations")

	if err != nil {
		fmt.Println(err)
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

func execMigrateUp(w http.ResponseWriter, r *http.Request) {

	addresses := r.Header.Get("X-Forwarded-For")
	fmt.Println(addresses)

	remote := r.RemoteAddr
	fmt.Println(remote)

	records, err := execMigrate(migrate.Up)
	if err != nil {
		fmt.Println(err)
	}

	bytes, _ := json.Marshal(records)

	recordString := string(bytes)

	fmt.Fprintf(w, "%s\n", recordString)

	// fmt.Printf("Applied %d migrations!\n", n)
}

func execMigrateDown(w http.ResponseWriter, r *http.Request) {

	addresses := r.Header.Get("X-Forwarded-For")
	fmt.Println(addresses)

	remote := r.RemoteAddr
	fmt.Println(remote)

	records, err := execMigrate(migrate.Down)
	if err != nil {
		fmt.Println(err)
	}

	bytes, _ := json.Marshal(records)

	recordString := string(bytes)

	fmt.Fprintf(w, "%s\n", recordString)

	// fmt.Printf("Applied %d migrations!\n", n)
}
