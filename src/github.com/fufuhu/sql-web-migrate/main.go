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
	http.HandleFunc("/migrate", execMigrate)

	// ListenするIPアドレスを定義
	err := http.ListenAndServe("0.0.0.0:18080", nil)

	if err != nil {
		fmt.Println(err)
	}

}

type LogRecord struct {
	ID        string `json:id`
	AppliedAt string `json:appliedAt`
}

func execMigrate(w http.ResponseWriter, r *http.Request) {

	addresses := r.Header.Get("X-Forwarded-For")
	fmt.Println(addresses)

	remote := r.RemoteAddr
	fmt.Println(remote)

	sourcePath := config.GetMigrationSourcePath()
	fmt.Println(sourcePath)

	source := migrate.FileMigrationSource{
		Dir: sourcePath,
	}

	// port, _ := config.GetPort()
	// sslMode, _ := config.GetSSLMode()

	//connectionString := "host=127.0.0.1 port=5432 user=migrate password=migrate dbname=migrate sslmode=disable"
	// connectionString := config.BuildConnectionString(
	// 	config.GetHost(),
	// 	port,
	// 	config.GetUser(),
	// 	config.GetPassword(),
	// 	config.GetDBName(),
	// 	sslMode)

	//https://godoc.org/github.com/lib/pq

	// connectionString := "postgres://migrate:migrate@migrate/?host=/var/run/postgresql"
	// connectionString := "host=/var/run/postgresql user=migrate password=migrate dbname=migrate"
	connectionString := config.BuildConnectionStringForUnixDomainSocket(
		config.GetHost(),
		config.GetUser(),
		config.GetPassword(),
		config.GetDBName())

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Println("DB connection open failure.")
		fmt.Println(err)
	}

	n, err := migrate.ExecMax(db, "postgres", source, migrate.Up, 0)

	if err != nil {
		fmt.Println(err)
	}

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

	records := []LogRecord{}
	for rows.Next() {
		err := rows.Scan(&id, &appliedAt)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, appliedAt)
		record := LogRecord{
			ID:        id,
			AppliedAt: appliedAt.String(),
		}

		records = append(records, record)

		// fmt.Fprintf(w, "%s %s\n", id, appliedAt)

	}

	bytes, _ := json.Marshal(records)

	recordString := string(bytes)

	fmt.Fprintf(w, "%s\n", recordString)

	fmt.Printf("Applied %d migrations!\n", n)
}
