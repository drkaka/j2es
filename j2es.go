package j2es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/drkaka/j2es/fetcher"
	"github.com/drkaka/j2es/uploader"
	"github.com/drkaka/lg"
	client "github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

var (
	// to record service cursor location
	records map[string]string
	// to record how many logs uploaded
	uploads map[string]int
)

func writeRecords() error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configInfo.Path, b, 0644); err != nil {
		return err
	}
	return nil
}

// get "journalctl" command and its arguments
func getCommand(serviceName, cursor string) (string, []string) {
	args := []string{"-u", serviceName, "-o", "json"}

	if len(cursor) != 0 {
		// if cursor existed, then find logs after the cursor
		args = append(args, fmt.Sprintf("--after-cursor=%s", cursor))
	}
	return "journalctl", args
}

func handleOne(serviceName string) error {
	uploads[serviceName] = 0
	defer lg.L(nil).Debug(serviceName, zap.Int("count", uploads[serviceName]))

	// setup command arguments
	cmd, args := getCommand(serviceName, records[serviceName])
	lg.L(nil).Debug(serviceName, zap.String("args", strings.Join(args, " ")))

	results, err := fetcher.GetMessages(serviceName, cmd, args...)
	l := len(results)

	if err != nil {
		return err
	} else if l == 0 {
		return nil
	}

	// add result to batch
	for _, one := range results {
		uploader.AddRecord(serviceName, one.Message)
	}

	// post to es hosts
	if err := uploader.PostToES(); err != nil {
		return err
	}

	// record the last cursor
	records[serviceName] = results[l-1].Cursor
	if err := writeRecords(); err != nil {
		return err
	}

	// record the log count
	uploads[serviceName] = l

	return nil
}

// recordStatus to influxdb
func recordStatus(started int64, success bool) error {
	if configInfo.IFInfo.Host == "" {
		// no need to upload to influxdb
		return nil
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     configInfo.IFInfo.Host,
		Username: configInfo.IFInfo.Username,
		Password: configInfo.IFInfo.Password,
	})
	if err != nil {
		return err
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  configInfo.IFInfo.Database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	// Create a point and add to batch
	fields := make(map[string]interface{})
	// add the uploads count
	for k, v := range uploads {
		fields[k] = v
	}
	// add duration
	fields["duration"] = time.Now().UnixNano()/1e6 - started
	// add status
	fields["status"] = success

	pt, err := client.NewPoint("j2es", nil, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		return err
	}
	return nil
}

// Start j2es with the config file.
func Start(configFile string) error {
	// ms
	started := time.Now().UnixNano() / 1e6

	if err := makeConfig(configFile); err != nil {
		lg.L(nil).Error("error make config", zap.Error(err))
		if err1 := recordStatus(started, false); err1 != nil {
			lg.L(nil).Error("error recording", zap.Error(err1))
		}
		return err
	}

	success := true
	for _, one := range configInfo.Services {
		if err := handleOne(one); err != nil {
			success = false
			lg.L(nil).Error("error handle service", zap.Error(err), zap.String("service", one))
		}
	}

	if err := recordStatus(started, success); err != nil {
		lg.L(nil).Error("error recording", zap.Error(err))
	}

	return nil
}
