package tackdb

import (
	"io/ioutil"
	"testing"
)

func TestConfigMerge(t *testing.T) {
	config := NewDefaults()
	data := []byte(`{"port":"75000"}`)
	err := config.merge(data)
	if err != nil {
		t.Errorf("Expected %s to be nil", err)
	}
	if config.Port != "75000" {
		t.Errorf("Expected %s to overwrite port field, but got %q", data, config)
	}

	data, err = ioutil.ReadFile("./testdata/invalidconfig.json")
	if err != nil {
		t.Fatal(err)
	}
	err = config.merge(data)
	if err == nil {
		t.Errorf("Expected %s to not be nil", err)
	}
}

func TestInitConfig(t *testing.T) {
	err := InitConfig("./testdata/invalidconfig.json")
	if err == nil {
		t.Errorf("Expected %s to not be nil", err)
	}
}
