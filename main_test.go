package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := loadConfig("./config.json")
	if err != nil {
		t.Fail()
	}
	if c.Cert == "" {
		t.Fail()
	}
}

func TestLoadConfigNotExisting(t *testing.T) {
	_, err := loadConfig("/tmp/rand")
	if !strings.Contains(err.Error(), "no such file") {
		t.Fail()
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	const filename = "/tmp/invalidJSON"
	ioutil.WriteFile(filename, []byte("invalid data"), 0777)

	_, err := loadConfig(filename)
	// we expect error to be of *json.SyntaxError type
	if _, expectedErrType := err.(*json.SyntaxError); !expectedErrType {
		t.Fail()
	}

	defer os.Remove(filename)
}
