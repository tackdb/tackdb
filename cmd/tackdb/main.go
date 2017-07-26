package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	flag "github.com/ogier/pflag"

	"gitlab.com/tackdb/tackdb"
)

var cli = flag.Bool("cli", false, "Connect as client")
var addr = flag.String("addr", ":3750", "Connection address.")

func main() {
	flag.Parse()

	if *cli {
		log.Fatal(runClient())
	} else {
		log.Fatal(tackdb.Serve())
	}
}

func runClient() error {
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		return err
	}

	serverout := bufio.NewReader(conn)
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := stdin.ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Fprintf(conn, input)

		msg, err := serverout.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Server disconnected.")
			return err
		} else if err != nil {
			return err
		}

		msg = strings.Trim(msg, "\n")
		fmt.Println(msg)
	}
}
