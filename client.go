package tackdb

import "net"

// Maintains a connection to the client.
type Client struct {
	conn *net.Conn
}
