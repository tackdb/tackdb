package tackdb

import (
	"bufio"
	"errors"
	"log"
	"net"
	"path"
	"sync"
)

var ErrMaxConn = errors.New("Reached max connections.")

type Server struct {
	clientid int64
	listener net.Listener
	handler  func() (net.Conn, error)
	store    *store
	mu       sync.RWMutex
}

func NewServer() *Server {
	fp := path.Join(*configdir, *configname)
	if err := ReadConfig(fp); err != nil {
		log.Printf("Error reading config file (%q): %s", fp, err)
	}

	s := &Server{
		store: NewStore(),
		mu:    sync.RWMutex{},
	}
	s.handler = s.Accept
	return s
}

func (s *Server) Accept() (net.Conn, error) {
	return s.listener.Accept()
}

func (s *Server) NoAccept() (net.Conn, error) {
	return nil, ErrMaxConn
}

func (s *Server) Listen() *Server {
	var err error
	log.Println("Listening on", config.Port)
	s.listener, err = net.Listen(SCHEME, ":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *Server) Serve() error {
	defer s.listener.Close()
	for {
		s.clientid++
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Make a new client and serve in a new goroutine.
		go s.NewClient(conn).Handle()
	}
}

func (s *Server) NewClient(conn net.Conn) *client {
	c := &client{
		conn:   conn,
		id:     s.clientid,
		reader: bufio.NewReader(conn),
	}
	c.commands = c.NewCommandTable(s)
	s.clientid++
	return c
}
