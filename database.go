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

// Setup configures the database connection's config information and creates the table
func Setup() {
	createDatabase()
	setDatabaseInfo()
	createTable()
}

// HasHash checks if the passed in hash is currently in the database.
func HasHash(hash string) bool {
	sqlStatement := fmt.Sprintf(`SELECT COUNT(*) FROM %v WHERE hash=$1;`, Config.Database.DBname)
	return checkValue(sqlStatement, hash)
}

// HasName checks if the passed in new name is used in the database.
func HasName(name string) bool {
	sqlStatement := fmt.Sprintf(`SELECT COUNT(*) FROM %v WHERE n_name=$1;`, Config.Database.DBname)
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

// createDatabase creates the database if it does not already exist
func createDatabase() {
	qry := fmt.Sprintf(`CREATE DATABASE %v WITH OWNER %v`, Config.Database.DBname, Config.Database.User)

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
			Config.Database.Host, Config.Database.Port, Config.Database.User,
			Config.Database.Password))

	if err != nil {
		log.Println("Error connecting to the db")
		log.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec(qry)
	if err != nil {
		log.Println(err)
		log.Println("Error executing create database command")
	}
	return
}

// createTable creates the table if it doesn't exist
func createTable() {
	qry := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %v (
		hash text PRIMARY KEY,
		o_name text,
		n_name text,
		date date,
		size integer
	)`, Config.Database.TableName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Error connecting to the db")
		log.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec(qry)
	if err != nil {
		log.Println(err)
		log.Println("Error executing create table command")
	}
	return
}

// setDatabaseInfo sets the Database information from the Config file
func setDatabaseInfo() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Config.Database.Host, Config.Database.Port, Config.Database.User,
		Config.Database.Password, Config.Database.DBname)
}
