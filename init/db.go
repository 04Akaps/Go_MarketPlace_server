package init

import (
	"database/sql"
	"fmt"
	sqlc "goServer/mysql/sqlc"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDBClient(driver, username, password, table, uri, port string) *sqlc.Queries {
	sourceUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, uri, port, table)
	dbInstance, err := sql.Open(driver, sourceUri)
	// RDS만들어서 Connect진행

	if err != nil {
		log.Fatal(err)
	}

	dbInstance.SetConnMaxLifetime(time.Minute * 1)
	dbInstance.SetMaxIdleConns(3)
	dbInstance.SetMaxOpenConns(6)

	return sqlc.New(dbInstance)
}
