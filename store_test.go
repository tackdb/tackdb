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
	"fmt"
	"strings"
	"testing"
)

var testcases = []struct {
	name     string
	in, want string
}{
	{"test 1", `SET ex 10
GET ex
UNSET ex
GET ex
END`, `
10

NULL
`}, {"test 2", `SET a 10
SET b 10
NUMEQUALTO 10
NUMEQUALTO 20
SET b 30
NUMEQUALTO 10
END`, `

2
0

1
`}, {"test 3", `BEGIN
SET a 10
GET a
BEGIN
SET a 20
GET a
ROLLBACK
GET a
ROLLBACK
GET a
END`, `

10


20

10

NULL
`}, {"test 4", `BEGIN
SET a 30
BEGIN
SET a 40
COMMIT
GET a
ROLLBACK
COMMIT
END`, `




40
NO TRANSACTION
NO TRANSACTION
`}, {"test 5", `SET a 50
BEGIN
GET a
SET a 60
BEGIN
UNSET a
GET a
ROLLBACK
GET a
COMMIT
GET a
END`, `

50



NULL

60

60
`}, {"test 6", `SET a 10
BEGIN
NUMEQUALTO 10
BEGIN
UNSET a
NUMEQUALTO 10
ROLLBACK
NUMEQUALTO 10
COMMIT
END`, `

1


0

1

`}, {"test 7", `SET foo bar baz
NUMEQUALTO bar baz
GET foo`, `
1
bar baz`},
}

func TestStore(t *testing.T) {
	for _, test := range testcases {
		store := NewStore()
		buffer := make([]string, 0)
		args := strings.Split(test.in, "\n")

		for _, argstr := range args {
			cmds := strings.Split(argstr, " ")
			resp, err := store.commands[cmds[0]](cmds[1:]...)
			if err != nil {
				resp = fmt.Sprintf("%s", err)
			}
			buffer = append(buffer, resp)
		}

		got := strings.Join(buffer, "\n")

		if got != test.want {
			t.Errorf("%q: %q != %q", test.name, got, test.want)
		}
	}
}
