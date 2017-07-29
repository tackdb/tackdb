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
	"errors"
	"fmt"
	"strings"
)

// command describes the form that parsed arguments must conform to.
type Command func(...string) (string, error)

// table is the basic key-value datatable type.
type table map[string]*string

// counter maps values to their counts.
type counter map[string]int

// store manages the key-value tables.
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
