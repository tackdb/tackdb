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
	"encoding/json"
	"io/ioutil"
	"os"

	flag "github.com/ogier/pflag"
)

const (
	SCHEME  = "tcp"
	VERSION = "0.1.0"
)

var (
	configname = flag.String("confname", "tackdb.conf", "Filename of TackDB runtime configuration file.")
	configdir  = flag.StringP("dir", "d", os.Getenv("HOME"), "Directory location of runtime configuration file (.tackrc).")
	port       = flag.StringP("port", "p", "3750", "TCP service port.")
	// maxconns   = flag.IntP("max-connections", "m", 0, "Maximum connections, setting to 0 will not limit the number of connections.")
	// adminname  = flag.String("admin-name", "admin", "Username of admin user.")
	// adminpass  = flag.String("admin-pass", "pass", "Password of admin user.")
)

type Config struct {
	Port           string `json:"port"`
	MaxConnections int    `json:"max_connections"`
	AdminName      string `json:"admin_username"`
	AdminPass      string `json:"admin_password"`
}

// Set configuration to defaults.
var config Config

func NewDefaults() *Config {
	return &Config{
		Port: *port,
		// MaxConnections: *maxconns,
		// AdminName:      *adminname,
		// AdminPass:      *adminpass,
	}
}

func ReadConfig(path string) (err error) {
	if data, err := ioutil.ReadFile(path); err == nil {
		return config.merge(data)
	}
	return
}

func (c *Config) merge(data []byte) (err error) {
	err = json.Unmarshal(data, c)
	return
}
