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
