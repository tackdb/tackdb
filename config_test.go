package tackdb

import (
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
}

func TestReadConfig(t *testing.T) {
	err := ReadConfig("./testdata/invalidconfig.json")
	if err == nil {
		t.Errorf("Expected %s to not be nil", err)
	}

	err = ReadConfig("./testdata/validconfig.json")
	if err != nil {
		t.Errorf("Expected %s to be nil", err)
	}
}
