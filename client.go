package tackdb

import (
	"bufio"
	// "container/list"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

// type Client interface {
// 	GetCommandTable() map[string]Command
// }

// client should peek first line
// if the first line is an AUTH statement,
// authenticate client and then connect
// if first line is not authenticate,
// try to add to pool
// Maintains a connection to the client.
type client struct {
	conn   net.Conn
	id     int64
	reader *bufio.Reader
	// queue    list.List
	commands map[string]Command
	hasLock  bool
}

func (c *client) Handle() {
	defer c.conn.Close()
	var err error

	log.Println(c.id)

	for {
		if err = c.ReadCommand(); err != nil {
			break
		}
	}
}

func (c *client) ReadCommand() error {
	string, err := c.reader.ReadString('\n')
	if err == io.EOF {
		if c.hasLock {
			if rollback, ok := c.commands["ROLLBACK"]; ok {
				var err error
				for ; err == nil; _, err = rollback() {
				}
			}
			// defer c.Unlock()
			// defer s.mu.Unlock()
		}
		return err
	}
	log.Println(err, string)

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

var ErrUnrecognizedCommand = errors.New("UNRECOGNIZED COMMAND")

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

var ErrNoLock = errors.New("NOT LOCKED")

func (c *client) NewCommandTable(s *Server) map[string]Command {
	table := make(map[string]Command)

	table["GET"] = func(args ...string) (string, error) {
		return RLockUnlock(c, s, func() (string, error) {
			return s.store.get(args[1:]...)
		})
	}
	table["SET"] = func(args ...string) (string, error) {
		return LockUnlock(c, s, func() (string, error) {
			if ret, err := s.store.stash(args[1:]...); err != nil {
				return ret, err
			}
			return s.store.set(args[1:]...)
		})
	}
	table["UNSET"] = func(args ...string) (string, error) {
		return LockUnlock(c, s, func() (string, error) {
			if ret, err := s.store.stash(args[1:]...); err != nil {
				return ret, err
			}
			return s.store.unset(args[1:]...)
		})
	}
	table["NUMEQUALTO"] = func(args ...string) (string, error) {
		return RLockUnlock(c, s, func() (string, error) {
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

func LockUnlock(c *client, s *Server, cb func() (string, error)) (string, error) {
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

func RLockUnlock(c *client, s *Server, cb func() (string, error)) (string, error) {
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
