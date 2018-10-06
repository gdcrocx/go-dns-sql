package godnssql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type dnslookup struct {
	dn string
	ip string
}

var wg = sync.WaitGroup{}

// CreateDBConnection : Returns a MySQL DB Connection
func CreateDBConnection() (db *sql.DB) {

	db, err := sql.Open("mysql", "golang:golang@tcp(127.0.0.1:3306)/golang_db")

	// log.Printf(reflect.TypeOf(db).String())

	if err != nil {
		panic(err.Error())
	}

	return db
}

// SelectDNSLookup : Returns a DNS Record from MySQL
func SelectDNSLookup(db *sql.DB, dn string) string {

	wg.Add(1)

	var results dnslookup

	// fmt.Println("SelectDNSLookup: " + dn)
	// fmt.Println("SelectDNSLookup: " + reflect.TypeOf(dn))

	err := db.QueryRow("SELECT dn, ip FROM golang_db.dnslookup WHERE dn = ?", string(dn)).Scan(&results.dn, &results.ip)

	if err != nil {
		log.Printf(err.Error())
		// fmt.Println("SelectDNSLookup: No record found.")
		// panic(err.Error())
		wg.Done()
		return "No record found"
	}

	wg.Done()
	fmt.Println("SelectDNSLookup: Found ", results.dn)
	return results.ip

}

// InsertDNSAlias : Insert a new DNS record in MySQL DB
func InsertDNSAlias(db *sql.DB, dn string, ip string) {

	wg.Add(1)

	stmtIns, err := db.Query("INSERT INTO golang_db.dnslookup (dn, ip) VALUES (?,?);", dn, ip)

	if err != nil {
		log.Printf(err.Error())
		// panic(err.Error())
		recover()
	} else {
		fmt.Println("InsertDNSAlias : ", dn, " = ", ip, " Added.")
	}

	defer stmtIns.Close()
	wg.Done()

}

func updateDNSAlias(db *sql.DB, dn string, ip string) {

	wg.Add(1)

	var results dnslookup

	// var buffer bytes.Buffer
	// buffer.WriteString("%")
	// buffer.WriteString(strings.TrimSpace(dn))
	// buffer.WriteString("%")

	err := db.QueryRow("SELECT dn, ip FROM golang_db.dnslookup WHERE dn = ?;", strings.TrimSpace(dn)).Scan(&results.dn, &results.ip)

	if err != nil {
		log.Println("Error : ", err.Error())
		if strings.Contains(err.Error(), "no rows in result set") {
			stmtIns, err := db.Query("INSERT INTO golang_db.dnslookup (dn, ip) VALUES (?,?);", dn, ip)

			if err != nil {
				log.Println("Error : ", err.Error())
				// panic(err.Error())
				recover()
			} else {
				log.Println("Added", dn, " = ", ip, ".")
			}

			defer stmtIns.Close()
		}
		recover()
	} else {
		if results.dn != "" && results.ip != "" {
			log.Println("Found", results.dn, results.ip)
			stmtUpd, err := db.Query("UPDATE golang_db.dnslookup SET ip = ? WHERE dn = ?;", ip, dn)

			if err != nil {
				log.Println("Error : ", err.Error())
				recover()
			} else {
				log.Println("Updated", dn, " = ", ip, ".")
			}

			defer stmtUpd.Close()
		}
	}
	wg.Done()

}

func dropDNSAlias(db *sql.DB, dn string) {

	wg.Add(1)

	// var buffer bytes.Buffer
	// buffer.WriteString("%")
	// buffer.WriteString(strings.Trim(dn, "\n"))
	// buffer.WriteString("%")

	err := db.QueryRow("DELETE * FROM golang_db.dnslookup WHERE dn = ?;", strings.TrimSpace(dn))

	if err != nil {
		log.Println(err)
		recover()
	} else {
		log.Println("Deleted", dn)
	}

	wg.Done()

}
