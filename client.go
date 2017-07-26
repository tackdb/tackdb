package tackdb

import (
	"bufio"
	"container/list"
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
	conn     net.Conn
	id       int64
	reader   *bufio.Reader
	queue    list.List
	commands map[string]Command
}

func (c *client) Run() {
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
		return err
	}
	log.Println(err, string)

	string = strings.TrimSpace(string)
	args := strings.Split(string, " ")

	if cmd, err := c.GetCommand(args); err != nil {
		c.conn.Write([]byte(err.Error() + "\n"))
	} else {
		resp, err := cmd(args[1:]...)
		if err != nil {
			c.conn.Write([]byte(err.Error() + "\n"))
		} else {
			c.conn.Write([]byte(resp + "\n"))
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
	name := strings.ToUpper(args[0])
	cmd, ok := c.commands[name]
	if !ok {
		return nil, ErrUnrecognizedCommand
	}
	return cmd, nil
}

// type User struct {
// 	conn *net.Conn
// }
//
// func (u *User) GetCommandTable() map[string]Command {
// 	return nil
// }
//
// type Admin struct {
// }
//
// func (a *Admin) GetCommandTable() map[string]Command {
// 	return nil
// }
