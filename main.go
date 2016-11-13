package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"

	"io/ioutil"

	"fmt"
	"j2es/fetcher"
	"j2es/logger"
	"j2es/uploader"
	"strings"

	"github.com/uber-go/zap"
)

type config struct {
	EsHosts []string `json:"es_hosts"`
	Names   []string `json:"services"`
	Path    string   `json:"record_path"`
}

var (
	configInfo config
	records    map[string]string
)

func makeConfig(path string) error {
	// read config file to configInfo.
	if b, err := ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("Error read config file: %s, %v", path, err)
	} else if err := json.Unmarshal(b, &configInfo); err != nil {
		return err
	}

	if len(configInfo.EsHosts) == 0 {
		return errors.New("\"es_hosts\" field is not existed")
	}

	if len(configInfo.Names) == 0 {
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
		// not exist
		records = make(map[string]string, 0)
	}

	return nil
}

func writeRecords(records map[string]string) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configInfo.Path, b, 0644); err != nil {
		return err
	}
	return nil
}

func getCommand(serviceName, cursor string) (string, []string) {
	var args []string
	args = append(args, "-u")
	args = append(args, serviceName)
	args = append(args, "-o")
	args = append(args, "json")

	if len(cursor) != 0 {
		args = append(args, fmt.Sprintf("--after-cursor=%s", cursor))
	}
	return "journalctl", args
}

func handleOne(serviceName string) error {
	cmd, args := getCommand(serviceName, records[serviceName])

	logger.Log.Debug(serviceName, zap.String("cmd", strings.Join(args, " ")))
	results, err := fetcher.GetMessages(serviceName, cmd, args...)
	l := len(results)
	if err != nil {
		return err
	} else if l == 0 {
		logger.Log.Info("", zap.String("args", strings.Join(args, " ")), zap.Int("count", l))
		return nil
	}

	for _, one := range results {
		uploader.AddRecord(serviceName, one.Message)
	}

	if err := uploader.PostToES(); err != nil {
		return err
	}

	records[serviceName] = results[l-1].Cursor
	if err := writeRecords(records); err != nil {
		return err
	}
	logger.Log.Info("", zap.String("args", strings.Join(args, " ")), zap.Int("count", l))

	return nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "path", "", "The config file path.")
	flag.Parse()

	if len(configPath) == 0 {
		logger.Log.Error("path is empty.")
		os.Exit(1)
	}

	if err := makeConfig(configPath); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	if err := uploader.PrepareUploader(configInfo.EsHosts); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	exCode := 0
	for _, one := range configInfo.Names {
		if err := handleOne(one); err != nil {
			logger.Log.Error(err.Error())
			exCode = 1
		}
	}
	os.Exit(exCode)
}
