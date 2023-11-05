package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// func InitializeDB() {
// 	mysqlUser := "test_user"
// 	mysqlPwd := "password"
// 	mysqlDatabase := "test_hackathon"

// 	connStr := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", mysqlUser, mysqlPwd, "localhost", mysqlDatabase)
// 	var err error
// 	DB, err = sql.Open("mysql", connStr)
// 	if err != nil {
// 		log.Fatal("Error opening database:", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		log.Fatalf("Error pinging the database: %v", err)
// 	}

// 	fmt.Printf("Database connection successful.")
// }

func InitializeDB() {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep passwords and other secrets safe.
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
		return
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName, opts...)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Errorf("sql.Open: %w", err)
		return
	}
	DB = dbPool
}
