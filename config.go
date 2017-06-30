package j2es

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	configInfo config
)

type influxdb struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"db"`
}

type config struct {
	EsHosts  []string `yaml:"es_hosts"`
	Services []string `yaml:"services"`
	Path     string   `yaml:"record_path"`
	IFInfo   influxdb `yaml:"influxdb,omitempty"` // if none, no post
}

func makeConfig(path string) error {
	// read config file to configInfo.
	if b, err := ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("Error read config file: %s, %v", path, err)
	} else if err := yaml.Unmarshal(b, &configInfo); err != nil {
		return err
	}

	if len(configInfo.EsHosts) == 0 {
		return errors.New("\"es_hosts\" field is not existed")
	}

	if len(configInfo.Services) == 0 {
		return errors.New("\"services\" field is not existed")
	}

	if len(configInfo.Path) == 0 {
		return errors.New("\"record_path\" field is not existed")
	}

	// read records file to services.
	if b, err := ioutil.ReadFile(configInfo.Path); err == nil {
		if err := json.Unmarshal(b, &records); err != nil {
			return err
		}

	} else {
		// no records found
		records = make(map[string]string, 0)
	}

	// init uploads
	uploads = make(map[string]int)
	return nil
}
