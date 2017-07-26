package tackdb

import (
	"bufio"
	"errors"
	"log"
	"net"
	"path"
	"sync"
)

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

var ErrMaxConn = errors.New("Reached max connections.")

func (s *Server) NoAccept() (net.Conn, error) {
	return nil, ErrMaxConn
}

func (s *Server) Listen() *Server {
	var err error
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

		// Make a new client.
		newcommand := make(map[string]Command)
		for name, cmd := range s.store.commands {
			newcommand[name] = cmd
		}

		client := &client{
			conn:   conn,
			id:     s.clientid,
			reader: bufio.NewReader(conn),
			// commands: newcommand,
		}
		client.addcommandlock(newcommand, &s.mu)
		go client.Run()
	}
}

func (c *client) addcommandlock(commands map[string]Command, mu *sync.RWMutex) {
	ogbegin := commands["BEGIN"]
	commands["BEGIN"] = func(args ...string) (string, error) {
		c.isLocked = true
		mu.Lock()
		return ogbegin(args...)
	}

	ogset := commands["SET"]
	commands["SET"] = func(args ...string) (string, error) {
		if c.isLocked {
			return ogset(args...)
		} else {
			mu.Lock()
			defer mu.Unlock()
			return ogset(args...)
		}
	}
}
