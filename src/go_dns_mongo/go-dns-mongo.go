package godnsmongo

import (
	"github.com/globalsign/mgo"
)

// GetSession : Returns a MongoDB Session to connect with MongoDB Server(s)
func GetSession() {
	Host := []string{
		"127.0.0.1:27017",
		// replica set addrs...
	}

	const (
		Username   = "YOUR_USERNAME"
		Password   = "YOUR_PASS"
		Database   = "YOUR_DB"
		Collection = "YOUR_COLLECTION"
	)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
		// Username: Username,
		// Password: Password,
		// Database: Database,
		// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
		// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
		// },
	})

	if err != nil {
		panic(err)
	}

	defer session.Close()

	return session
}
