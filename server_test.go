package tackdb

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer().Listen()
	go server.Serve()
	defer server.listener.Close()

	client, err := net.Dial("tcp", ":"+*port)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	response := bufio.NewReader(client)

	fmt.Fprintf(client, "GET a\n")

	msg, err := response.ReadString('\n')
	if msg != "NULL\n" {
		t.Errorf("Expected %s to be NULL", msg)
	}
}
