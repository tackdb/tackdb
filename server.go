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

// Server manages the TCP connection and client pool.
// Each client references the data store and mutex.
type Server struct {
	clientid int64
	listener net.Listener
	handler  func() (net.Conn, error)
	store    *store
	mu       sync.RWMutex
}

func NewServer() *Server {
	config = *NewDefaults()

	fp := path.Join(*configdir, *configname)
	if err := ReadConfig(fp); err != nil {
		log.Printf("Error reading config file (%q): %s", fp, err)
	}

	s := &Server{
		store:    NewStore(),
		mu:       sync.RWMutex{},
		clientid: 1,
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
	log.Println("TackDB", "v"+VERSION, "Listening on", config.Port)
	s.listener, err = net.Listen(SCHEME, ":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *Server) Serve() error {
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
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
	c.commands = c.newCommandTable(s)
	s.clientid++
	return c
}
