package main

import (
	"log"

	flag "github.com/ogier/pflag"
	"github.com/tackdb/tackdb"
)

func main() {
	flag.Parse()
	log.Fatal(tackdb.NewServer().Listen().Serve())
}
