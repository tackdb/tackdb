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
	"fmt"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer().Listen()
	go server.Serve()
	defer server.listener.Close()

	client, err := net.Dial("tcp", ":"+*port)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	response := bufio.NewReader(client)

	fmt.Fprintf(client, "GET a\n")

	msg, err := response.ReadString('\n')
	if msg != "NULL\n" {
		t.Errorf("Expected %s to be NULL", msg)
	}
}
