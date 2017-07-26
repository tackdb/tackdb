package tackdb

import (
	"encoding/json"
	"io/ioutil"
	"os"

	flag "github.com/ogier/pflag"
)

const (
	SCHEME  = "tcp"
	VERSION = "0.0.1"
)

var (
	configname = flag.String("confname", "tackdb.conf", "Filename of TackDB runtime configuration file.")
	configdir  = flag.StringP("dir", "d", os.Getenv("HOME"), "Directory location of runtime configuration file (.tackrc).")
	port       = flag.StringP("port", "p", "3750", "Default port to bind to.")
	maxconns   = flag.IntP("max-connections", "m", 0, "Maximum connections, setting to 0 will not limit the number of connections.")
	adminname  = flag.String("admin-name", "admin", "Username of admin user.")
	adminpass  = flag.String("admin-pass", "pass", "Password of admin user.")
)

type Config struct {
	Port           string `json:"port"`
	MaxConnections int    `json:"max-connections"`
	AdminName      string `json:"admin_username"`
	AdminPass      string `json:"admin_password"`
}

// Set configuration to defaults.
var config = NewDefaults()

// temp to try IIFE
var temp = func() *Config {
	return &Config{
		Port:           *port,
		MaxConnections: *maxconns,
		AdminName:      *adminname,
		AdminPass:      *adminpass,
	}
}()

func NewDefaults() *Config {
	return &Config{
		Port:           *port,
		MaxConnections: *maxconns,
		AdminName:      *adminname,
		AdminPass:      *adminpass,
	}
}

func InitConfig(path string) (err error) {
	if data, err := ioutil.ReadFile(path); err == nil {
		return config.merge(data)
	}
	return
}

func (c *Config) merge(data []byte) (err error) {
	err = json.Unmarshal(data, c)
	return
}
