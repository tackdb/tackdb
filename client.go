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
	"io"
	"log"
	"net"
	"strings"
)

var ErrUnrecognizedCommand = errors.New("UNRECOGNIZED COMMAND")

// TODO: Basic authentication.
// client should peek first line
// if the first line is an AUTH statement,
// authenticate client and then connect
// if first line is not authenticate,
// return error if authentication is required

// client describes an individual connection to the server.
type client struct {
	conn     net.Conn
	id       int64
	reader   *bufio.Reader
	commands map[string]Command
	hasLock  bool
}

func (c *client) Handle() {
	defer c.conn.Close()
	var err error

	log.Println("Connection received. Client", c.id)

	for {
		if err = c.ReadCommand(); err != nil {
			break
		}
	}
}

func (c *client) ReadCommand() error {
	string, err := c.reader.ReadString('\n')
	if err == io.EOF {
		// Client closed the connection.
		// Cancel any outstanding transactions.
		if c.hasLock {
			if rollback, ok := c.commands["ROLLBACK"]; ok {
				var err error
				for ; err == nil; _, err = rollback() {
				}
			}
		}
		return err
	}

	string = strings.TrimSpace(string)
	args := strings.Split(string, " ")

	if cmd, err := c.GetCommand(args...); err != nil {
		c.conn.Write([]byte(err.Error() + "\n"))
	} else {
		resp, err := cmd(args...)
		if err != nil {
			c.conn.Write([]byte(err.Error() + "\n"))
		} else {
			c.conn.Write([]byte(resp + "\n"))
		}
	}
	return nil
}

func (c *client) GetCommand(args ...string) (Command, error) {
	if len(args) < 1 {
		return nil, ErrUnrecognizedCommand
	}
	name := strings.ToUpper(args[0])
	cmd, ok := c.commands[name]
	if !ok {
		return nil, ErrUnrecognizedCommand
	}
	return cmd, nil
}

func (c *client) Lock() {
	c.hasLock = true
}
func (c *client) Unlock() {
	c.hasLock = false
}

func (c *client) newCommandTable(s *Server) map[string]Command {
	table := make(map[string]Command)

	table["GET"] = func(args ...string) (string, error) {
		return rLockUnlock(c, s, func() (string, error) {
			return s.store.get(args[1:]...)
		})
	}
	table["SET"] = func(args ...string) (string, error) {
		return lockUnlock(c, s, func() (string, error) {
			if ret, err := s.store.stash(args[1:]...); err != nil {
				return ret, err
			}
			return s.store.set(args[1:]...)
		})
	}
	table["UNSET"] = func(args ...string) (string, error) {
		return lockUnlock(c, s, func() (string, error) {
			if ret, err := s.store.stash(args[1:]...); err != nil {
				return ret, err
			}
			return s.store.unset(args[1:]...)
		})
	}
	table["NUMEQUALTO"] = func(args ...string) (string, error) {
		return rLockUnlock(c, s, func() (string, error) {
			return s.store.numequalto(args[1:]...)
		})
	}
	table["BEGIN"] = func(...string) (string, error) {
		if c.hasLock {
			return s.store.begin()
		}
		s.mu.Lock()
		c.Lock()
		return s.store.begin()
	}
	table["COMMIT"] = func(...string) (string, error) {
		if !c.hasLock {
			return "", ErrNoTransaction
		}
		defer c.Unlock()
		defer s.mu.Unlock()
		return s.store.commit()
	}
	table["ROLLBACK"] = func(...string) (string, error) {
		if !c.hasLock {
			return "", ErrNoTransaction
		}
		// If we rollback the last transaction, unlock.
		if len(s.store.tables) == 2 {
			defer c.Unlock()
			defer s.mu.Unlock()
		}
		return s.store.rollback()
	}

	return table
}

// lockUnlock wraps the callback with the server mutex lock and client lock.
// If the client already has the lock, then execute the callback with locking.
// Though the server handles concurrent client connections, each client
// connection can only send commands sequentially.
func lockUnlock(c *client, s *Server, cb func() (string, error)) (string, error) {
	if c.hasLock {
		return cb()
	} else {
		s.mu.Lock()
		c.Lock()
		defer c.Unlock()
		defer s.mu.Unlock()
		return cb()
	}
}

// rLockUnlock wraps the callback with the server mutex read lock and client lock.
// If the client already has the lock, then execute the callback without read lock.
func rLockUnlock(c *client, s *Server, cb func() (string, error)) (string, error) {
	if c.hasLock {
		return cb()
	} else {
		s.mu.RLock()
		c.Lock()
		defer c.Unlock()
		defer s.mu.RUnlock()
		return cb()
	}
}
