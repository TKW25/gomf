package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type MetaData struct {
	hash   string
	o_name string
	n_name string
	size   int64
	date   time.Time
}

var psqlInfo string

// Setup configures the database connection's config information and creates the table
func Setup() {
	setDatabaseInfo()
}

// WriteMetadata writes the passed metadata to the database
func WriteMetadata(value MetaData) {
	statement := fmt.Sprintf(`
		INSERT INTO %v 
		("hash", "o_name", "n_name", "size", "date") 
		VALUES('%v', '%v', '%v', %d, '%v');
	`, Config.Database.TableName, value.hash, value.o_name, value.n_name, value.size, value.date.Format(time.ANSIC))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Error connecting to the db")
		log.Println(err)
		return
	}
	defer db.Close()

	if _, err = db.Query(statement); err != nil {
		log.Println("Error writing " + statement)
		log.Println(err)
		return
	}
}

// HasHash checks if the passed in hash is currently in the database.
func HasHash(hash string) bool {
	sqlStatement := fmt.Sprintf(`SELECT hash FROM %v WHERE hash=$1;`, Config.Database.TableName)
	return checkValue(sqlStatement, hash)
}

// HasName checks if the passed in new name is used in the database.
func HasName(name string) bool {
	sqlStatement := fmt.Sprintf(`SELECT n_name FROM %v WHERE n_name=$1;`, Config.Database.TableName)
	return checkValue(sqlStatement, name)
}

// // check runs the passed in statement checking if the value is in the database
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

// setDatabaseInfo sets the Database information from the Config file
func setDatabaseInfo() {
	psqlInfo = fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		Config.Database.Host, Config.Database.Port, Config.Database.User,
		Config.Database.Password, Config.Database.DBname)
}
