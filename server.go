package tackdb

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"path"

	flag "github.com/ogier/pflag"
)

type Server struct {
	clientid int64
	listener net.Listener
	handler  func() (net.Conn, error)
	store    *store
}

func NewServer() *Server {
	s := &Server{
		store: NewStore(),
	}
	s.handler = s.Accept
	return s
}

func (s *Server) Accept() (net.Conn, error) {
	return s.listener.Accept()
}

var ErrMaxConn = errors.New("Reached max connections.")

func (s *Server) NoAccept() (net.Conn, error) {
	return nil, ErrMaxConn
}

func (s *Server) Serve() error {
	for {
		s.clientid++
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Make a new client.
		newcommand := make(map[string]Command)
		for name, cmd := range s.store.commands {
			newcommand[name] = cmd
		}

		client := &client{
			conn:     conn,
			id:       s.clientid,
			reader:   bufio.NewReader(conn),
			commands: newcommand,
		}
		go client.Listen()
	}
}

func Serve() error {
	flag.Parse()
	fp := path.Join(*configdir, *configname)
	if err := InitConfig(fp); err != nil {
		log.Printf("Error reading config file (%q): %s", fp, err)
	}

	var err error
	server := NewServer()
	server.listener, err = net.Listen(SCHEME, ":"+config.Port)
	if err != nil {
		return err
	}
	defer server.listener.Close()

	//
	// server := NewServer()
	// server.listener = connection

	clientid := 1
	for {
		client, err := server.listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(client, clientid)
		}
		clientid++
	}
	return nil
}
