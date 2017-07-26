package tackdb

import (
	"fmt"
	"log"
	"net"
	"path"

	flag "github.com/ogier/pflag"
)

func Serve() error {
	flag.Parse()
	fp := path.Join(*configdir, *configname)
	if err := InitConfig(fp); err != nil {
		log.Printf("Error reading config file (%q): %s", fp, err)
	}

	connection, err := net.Listen(SCHEME, ":"+config.Port)
	if err != nil {
		return err
	}
	defer connection.Close()

	clientid := 1
	for {
		client, err := connection.Accept()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(client, clientid)
		}
		clientid++
	}
	return nil
}
