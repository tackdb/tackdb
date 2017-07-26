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

var cli = flag.Bool("cli", false, "Connect as client.")
var errcli = flag.Bool("errcli", false, "Run errant client.")
var addr = flag.String("addr", ":3750", "Connection address.")

func main() {
	flag.Parse()

	if *errcli {
		runErrantClient()
	} else if *cli {
		log.Fatal(runClient())
	} else {
		log.Fatal(tackdb.NewServer().Listen().Serve())
	}
}

func runErrantClient() {
	done := make(chan error)
	for i := 0; i < 20; i++ {
		go func() {
			conn, err := net.Dial("tcp", *addr)
			if err != nil {
				log.Println(err)
				done <- err
			}

			for j := 0; j < 20; j++ {
				fmt.Fprintf(conn, "SET foo bar\n")
			}
			done <- nil
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
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
