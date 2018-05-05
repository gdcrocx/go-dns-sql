package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type dnslookup struct {
	dn string
	ip string
}

var wg = sync.WaitGroup{}

func main() {

	wg.Add(1)
	go dnsListener(":53")

	wg.Add(1)
	go func() {
		log.Println("Hello World!")
		wg.Done()
	}()

	db := createDBConnection()

	defer db.Close()

	// err = db.Ping()

	// if err != nil {
	// 	panic(err.Error())
	// }

	// dn := "georgedavis.de"
	// ip := "127.0.0.1"

	// go selectDNSLookup(db, dn)
	// go insertDNSAlias(db, dn, ip)

	log.Println("Closing...")

	wg.Wait()

}

func createDBConnection() (db *sql.DB) {

	db, err := sql.Open("mysql", "golang:golang@tcp(127.0.0.1:3306)/golang_db")

	// log.Println(reflect.TypeOf(db).String())

	if err != nil {
		panic(err.Error())
	}

	return db
}

func selectDNSLookup(db *sql.DB, dn string) {

	wg.Add(1)

	var results dnslookup

	log.Println("Query Keyword - ", dn, reflect.TypeOf(dn).String())

	var buffer bytes.Buffer
	buffer.WriteString("%")
	buffer.WriteString(strings.Trim(dn, "\n"))
	buffer.WriteString("%")

	err := db.QueryRow("SELECT dn, ip FROM golang_db.dnslookup WHERE dn like ?;", buffer.String()).Scan(&results.dn, &results.ip)

	if err != nil {
		log.Println("Error : ", err.Error())
		// log.Println("No record found.")
		// panic(err.Error())
		recover()
	} else {
		log.Println("Found", results.dn, results.ip)
	}

	wg.Done()

}

func updateDNSAlias(db *sql.DB, dn string, ip string) {

	wg.Add(1)

	var results dnslookup

	var buffer bytes.Buffer
	buffer.WriteString("%")
	buffer.WriteString(strings.Trim(dn, "\n"))
	buffer.WriteString("%")

	err := db.QueryRow("SELECT dn, ip FROM golang_db.dnslookup WHERE dn like ?;", buffer.String()).Scan(&results.dn, &results.ip)

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
			stmtUpd, err := db.Query("UPDATE golang_db.dnslookup SET ip = ? WHERE dn like ?;", ip, dn)

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

	var buffer bytes.Buffer
	buffer.WriteString("%")
	buffer.WriteString(strings.Trim(dn, "\n"))
	buffer.WriteString("%")

	err := db.QueryRow("DELETE * FROM golang_db.dnslookup WHERE dn like ?;", buffer.String())

	if err != nil {
		log.Println(err)
		recover()
	} else {
		log.Println("Deleted", dn)
	}

	wg.Done()

}

func dnsListener(port string) {

	log.Println("Starting DNS Listener...")

	ln, err := net.Listen("tcp", ":53")

	if err != nil {
		// handle error
		log.Println("Error creating connection... Closing DNS Listener..")
		wg.Done()
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			// handle error
			log.Println("Error accepting connection...")
			wg.Done()
		}

		go handleConnection(conn)
	}

	log.Println("Signal for closure of DNS Listener...")

}

func handleConnection(conn net.Conn) {
	log.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		log.Println("Closing connection...")
		conn.Close()
	}()

	timeoutDuration := 3600 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read tokens delimited by newline
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}

		// log.Println("%s", bytes)
		// log.Println(reflect.TypeOf(bytes))
		inputString := strings.Trim(string(bytes[:len(bytes)]), "\n\t\r")

		db := createDBConnection()

		if strings.Contains(inputString, "update") {
			// Pulling keyword 'update' from the inputString
			inputString = strings.Replace(inputString, "update", "", 1)
			// Splitting strings with the '=' to get the domainName and ipAddress
			inputs := strings.Split(inputString, "=")
			go updateDNSAlias(db, inputs[0], inputs[1])
		} else if strings.Contains(inputString, "drop") {
			// Pulling keyword 'drop' from the inputString
			inputString = strings.Replace(inputString, "drop", "", 1)
			go dropDNSAlias(db, inputString)
		} else {
			go selectDNSLookup(db, inputString)
		}
	}
}
