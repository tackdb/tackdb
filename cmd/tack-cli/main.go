// Copyright 2017 Matthew Tso
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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
	"github.com/tackdb/tackdb"
)

var (
	cli    = flag.Bool("cli", false, "Connect as client.")
	errcli = flag.Bool("errcli", false, "Run errant client.")
	addr   = flag.String("addr", ":3750", "Connection address.")
)

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
	for i := 0; i < 200; i++ {
		go func() {
			conn, err := net.Dial("tcp", *addr)
			if err != nil {
				log.Println(err)
				done <- err
			}
			defer conn.Close()

			for j := 0; j < 200; j++ {
				fmt.Fprintf(conn, "GET foo\n")
				fmt.Fprintf(conn, "NUMEQUALTO bar\n")
				fmt.Fprintf(conn, "BEGIN\n")
				fmt.Fprintf(conn, "SET foo bar\n")
				fmt.Fprintf(conn, "UNSET foo\n")
				fmt.Fprintf(conn, "SET foo bar\n")
				fmt.Fprintf(conn, "ROLLBACK\n")
				fmt.Fprintf(conn, "GET foo\n")
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
	defer conn.Close()

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
