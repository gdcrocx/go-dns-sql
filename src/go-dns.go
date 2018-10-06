package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go_dns_redis"
	"go_dns_sql"

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

	// mongoDb := goDnsMongo.GetSession()

	// getDNSAlias("key3")
	// getDNSAlias("key45")
	// getDNSAlias("key2")

	// redisSetResponse := godnsredis.SetRedisKey(redisClient, "key3", "value4")

	// if strings.Contains(redisSetResponse, "SetRedisKeyError") {
	// 	panic("SetRedisKey returned error - " + redisSetResponse)
	// } else {
	// 	fmt.Println("SetRedisKey returned, Key: " + redisSetResponse)
	// }

	log.Printf("Closing...")

	wg.Wait()

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
			fmt.Println("DNS Listener: Error accepting connection...")
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
			fmt.Println("Connection Handler: " + err.Error())
		}

		// fmt.Printf("%s", bytes)
		// fmt.Println("Connection Handler: " + reflect.TypeOf(bytes))
		inputString := strings.TrimSpace(string(bytes[:len(bytes)]))

		if len(inputString) > 0 {
			go getDNSAlias(inputString)
		}
	}
}

func getDNSAlias(dn string) string {

	var dnsIP string

	redisClient := godnsredis.CreateRedisDatabaseConnection()

	isRedisKey := godnsredis.IsRedisKey(redisClient, dn)

	fmt.Println("Is Redis Key : " + strconv.FormatBool(isRedisKey))

	if isRedisKey {
		dnsIP = godnsredis.GetRedisKey(redisClient, dn)
	} else {

		sqlDb := godnssql.CreateDBConnection()

		defer sqlDb.Close()

		err := sqlDb.Ping()

		if err != nil {
			panic(err.Error())
		}

		// dn := "georgedavis.de"
		// ip := "127.0.0.1"

		// go godnssql.InsertDNSAlias(db, dn, ip)
		dnsIP = godnssql.SelectDNSLookup(sqlDb, dn)

		if !strings.Contains(dnsIP, "No record found") {
			godnsredis.SetRedisKey(redisClient, dn, dnsIP)
		}
	}

	fmt.Println("Get DNS Alias: " + dnsIP)
	return dnsIP
}
