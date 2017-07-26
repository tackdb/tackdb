package main

import (
	"log"

	"gitlab.com/tackdb/tackdb"
)

func main() {

	log.Fatal(tackdb.Serve())
}
