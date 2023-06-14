package db

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func InitDb() *sql.DB {
	connstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s")


}