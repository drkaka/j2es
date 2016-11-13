package main

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	if err := makeConfig("config_test/config_no_names.json"); err == nil {
		t.Logf("services: %v", configInfo.Names)
		t.Fatal("Should have error because of no services fields.")
	}

	configInfo.Names = []string{}
	configInfo.Path = ""

	if err := makeConfig("config_test/config_no_path.json"); err == nil {
		t.Logf("record_path: %s", configInfo.Path)
		t.Fatal("Should have error because of no record path fields.")
	}

	configInfo.Names = []string{}
	configInfo.Path = ""

	if err := makeConfig("config_test/config.json"); err != nil {
		t.Fatal(err)
	}

	serviceNames := []string{"a", "b"}
	if !reflect.DeepEqual(configInfo.Names, serviceNames) {
		t.Error("service name result wrong")
	}

	if records["a"] != "abc" {
		t.Error("records result wrong")
	}
}

func TestGetCommand(t *testing.T) {
	cmd, args := getCommand("abc", "cba")
	if cmd != "journalctl" {
		t.Error("command wrong.")
	}

	if !reflect.DeepEqual(args, []string{"-u", "abc", "-o", "json", "--after-cursor=cba"}) {
		t.Error("args wrong.", args)
	}
}
