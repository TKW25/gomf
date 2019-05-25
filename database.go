package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type RowValues struct {
	hash   string
	o_name string
	n_name string
	size   int
	date   time.Time
}

var psqlInfo string

// SetDatabaseInfo sets the Database information from the Config file
func SetDatabaseInfo() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Config.database.host, Config.database.port, Config.database.user,
		Config.database.password, Config.database.dbname)
}

// HasHash checks if the passed in hash is currently in the database.
func HasHash(hash string) bool {
	sqlStatement := fmt.Sprintf(`SELECT COUNT(*) FROM %v WHERE hash=$1;`, Config.database.dbname)
	return checkValue(sqlStatement, hash)
}

// HasName checks if the passed in new name is used in the database.
func HasName(name string) bool {
	sqlStatement := fmt.Sprintf(`SELECT COUNT(*) FROM %v WHERE n_name=$1;`, Config.database.dbname)
	return checkValue(sqlStatement, name)
}

// TODO: make type for row entries
func WriteMetadata(value RowValues) {

}

// check runs the passed in statement checking if the value is in the database
func checkValue(statement string, value string) bool {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Error connecting to the db")
		log.Println(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query(statement, value)
	if err != nil {
		log.Println("Error making a query")
		log.Println(err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}
