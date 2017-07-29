package tackdb

import (
	"errors"
	"fmt"
	"strings"
)

type Command func(...string) (string, error)

type Store interface {
	Get(string) (string, error)
	Set(string, string) (string, error)
	Unset(string) (string, error)
	Numequalto(string) (string, error)
	End() (string, error)
	Begin() (string, error)
	Rollback() (string, error)
	Commit() (string, error)
}

type table map[string]*string
type counter map[string]int

type store struct {
	tables   []table
	count    counter
	commands map[string]Command
}

func NewStore() (s *store) {
	s = &store{
		tables: []table{make(table)},
		count:  make(counter),
	}
	s.commands = map[string]Command{
		"GET": s.get,
		"SET": func(args ...string) (string, error) {
			s.stash(args[0])
			return s.set(args...)
		},
		"UNSET": func(args ...string) (string, error) {
			s.stash(args[0])
			return s.unset(args...)
		},
		"END":        identity,
		"NUMEQUALTO": s.numequalto,
		"BEGIN":      s.begin,
		"ROLLBACK":   s.rollback,
		"COMMIT":     s.commit,
	}
	return
}

func identity(...string) (string, error) {
	return "", nil
}

var ErrNil = errors.New("NULL")
var ErrNoTransaction = errors.New("NO TRANSACTION")

// Add a datatable to the end of the table slice.
func (s *store) begin(...string) (string, error) {
	s.tables = append(s.tables, make(table))
	return "", nil
}

func (s *store) stash(args ...string) (string, error) {
	if len(args) < 1 {
		return "", ErrInvalidArgs
	}
	key := args[0]
	// Check if we are in a transaction.
	if len(s.tables) < 2 {
		return "", nil
	}
	// If the value from the previous table has not been saved,
	// save it. If it has, don't override it.
	current := s.tables[len(s.tables)-1]
	// current := len(s.tables) - 1
	if _, ok := current[key]; !ok {
		if prev, ok := s.tables[0][key]; ok {
			current[key] = prev
		} else {
			current[key] = nil
		}
	}
	return "", nil
}

func (s *store) commit(...string) (string, error) {
	if len(s.tables) < 2 {
		return "", ErrNoTransaction
	}
	s.tables = s.tables[:1]
	return "", nil
}

// get retrieves the value for the given key in the first table of the
func (s *store) get(keys ...string) (string, error) {
	key := strings.Join(keys, " ")
	if value, ok := s.tables[0][key]; ok {
		return *value, nil
	} else {
		return "", ErrNil
	}
}

func (s *store) numequalto(values ...string) (string, error) {
	value := strings.Join(values, " ")
	count := s.count[value]
	return fmt.Sprintf("%d", count), nil
}

var ErrInvalidArgs = errors.New("INVALID ARGUMENTS")

func (s *store) set(args ...string) (string, error) {
	if len(args) < 2 {
		return "", ErrInvalidArgs
	}
	key, value := args[0], strings.Join(args[1:], " ")
	s.unset(key)
	s.tables[0][key] = &value
	s.count[value]++
	return "", nil
}

func (s *store) unset(keys ...string) (string, error) {
	key := strings.Join(keys, " ")
	value, err := s.get(key)
	if err == ErrNil {
		return "", nil
	}
	if s.count[value] > 1 {
		s.count[value] -= 1
	} else {
		delete(s.count, value)
	}
	delete(s.tables[0], key)
	return "", nil
}

func (s *store) rollback(...string) (string, error) {
	tablelen := len(s.tables)
	if tablelen < 2 {
		return "", ErrNoTransaction
	}

	var last table
	last, s.tables = s.tables[tablelen-1], s.tables[:tablelen-1]
	for key, value := range last {
		if value == nil {
			s.unset(key)
		} else {
			s.set(key, *value)
		}
	}
	return "", nil
}
