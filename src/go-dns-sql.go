package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"reflect"
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
		log.Printf("Hello World!")
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

	log.Printf("Closing...")

	wg.Wait()

}

func createDBConnection() (db *sql.DB) {

	db, err := sql.Open("mysql", "golang:golang@tcp(127.0.0.1:3306)/golang_db")

	// log.Printf(reflect.TypeOf(db).String())

	if err != nil {
		panic(err.Error())
	}

	return db
}

func selectDNSLookup(db *sql.DB, dn string) {

	wg.Add(1)

	var results dnslookup

	fmt.Println(dn)

	fmt.Println(reflect.TypeOf(dn))

	err := db.QueryRow("SELECT dn, ip FROM golang_db.dnslookup WHERE dn = ?", string(dn)).Scan(&results.dn, &results.ip)

	if err != nil {
		log.Printf(err.Error())
		// fmt.Println("No record found.")
		// panic(err.Error())
		recover()
	} else {
		fmt.Println("Found ", results.dn)
	}

	wg.Done()

}

func insertDNSAlias(db *sql.DB, dn string, ip string) {

	wg.Add(1)

	stmtIns, err := db.Query("INSERT INTO golang_db.dnslookup (dn, ip) VALUES (?,?);", dn, ip)

	if err != nil {
		log.Printf(err.Error())
		// panic(err.Error())
		recover()
	} else {
		fmt.Println(dn, " = ", ip, " Added.")
	}

	defer stmtIns.Close()
	wg.Done()

}

func dnsListener(port string) {

	fmt.Println("Starting DNS Listener...")

	ln, err := net.Listen("tcp", ":53")

	if err != nil {
		// handle error
		fmt.Println("Error creating connection... Closing DNS Listener..")
		wg.Done()
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			// handle error
			fmt.Println("Error accepting connection...")
			wg.Done()
		}

		go handleConnection(conn)
	}

	fmt.Println("Signal for closure of DNS Listener...")

}

func handleConnection(conn net.Conn) {
	fmt.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read tokens delimited by newline
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s", bytes)
		// fmt.Println(reflect.TypeOf(bytes))
		inputString := string(bytes[:len(bytes)])

		db := createDBConnection()

		go selectDNSLookup(db, inputString)
	}
}
