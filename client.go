package tackdb

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"net"
	"strings"
)

// type Client interface {
// 	GetCommandTable() map[string]Command
// }

// client should read first line
// if the first line is an AUTH statement,
// authenticate client and then connect
// if first line is not authenticate,
// try to add to pool
type client struct {
	conn     net.Conn
	id       int64
	reader   *bufio.Reader
	queue    list.List
	commands map[string]Command
}

func (c *client) Listen() {
	defer c.conn.Close()
	var err error

	for {
		if err = c.Accept(); err != nil {
			break
		}
	}
}

func (c *client) Accept() error {
	string, err := c.reader.ReadString('\n')
	if err == io.EOF {
		return err
	}

	string = strings.TrimSpace(string)
	args := strings.Split(string, " ")

	if cmd, err := c.GetCommand(args); err != nil {
		c.conn.Write([]byte(err.Error()))
	} else {
		resp, err := cmd(args[1:]...)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
		} else {
			c.conn.Write([]byte(resp))
		}
	}
	return nil
}

var ErrUnrecognizedCommand = errors.New("UNRECOGNIZED COMMAND")

// var ErrNoCommand = errors.New("NO COMMAND")
func (c *client) GetCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, ErrUnrecognizedCommand
	}
	cmd, ok := c.commands[args[0]]
	if !ok {
		return nil, ErrUnrecognizedCommand
	}
	return cmd, nil
}

// Maintains a connection to the client.
type User struct {
	conn *net.Conn
}

func (u *User) GetCommandTable() map[string]Command {
	return nil
}

type Admin struct {
}

func (a *Admin) GetCommandTable() map[string]Command {
	return nil
}
