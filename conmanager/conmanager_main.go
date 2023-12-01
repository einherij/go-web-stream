package conmanager

import (
	"fmt"
	"strconv"
	"time"
)

// Connection is a type for connecting to video stream, and than break it
type Connection struct {
	ConnectionID   int
	Timeout        time.Duration
	KeepAlive      chan bool
	MustDisconnect bool
}

var (
	// Connections is a dictionary of all connection on the server
	Connections  = make(map[string]*Connection)
	connectionid int // Incrementd connection id for unic value
	// timeout duration for connection
	timeoutDuration = time.Second * 20
)

// NewConnection creates new object of connections and returns it
func NewConnection() (c *Connection) {
	newCon := Connection{connectionid, timeoutDuration, make(chan bool, 1), false}
	Connections[strconv.Itoa(connectionid)] = &newCon
	fmt.Println("Start connection id:", newCon.ConnectionID)
	connectionid++
	return &newCon
}

// StartTimer begins timeout for connection
func (c *Connection) StartTimer() {
	select {
	case <-c.KeepAlive:
		c.StartTimer()
	case <-time.After(c.Timeout):
		c.MustDisconnect = true
	}
}

// KeepConnectionAlive sends keep-alive signal to update timer for connection
func KeepConnectionAlive(id string) {
	c := Connections[id]
	fmt.Println("Received keep alive request for connection id:", c.ConnectionID)
	c.KeepAlive <- true
}
