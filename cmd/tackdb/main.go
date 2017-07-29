package main

import (
	"log"

	flag "github.com/ogier/pflag"
	"gitlab.com/tackdb/tackdb"
)

func main() {
	flag.Parse()
	log.Fatal(tackdb.NewServer().Listen().Serve())
}
