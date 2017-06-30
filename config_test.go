package j2es

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeConfig(t *testing.T) {
	err := makeConfig("config_test/config_no_services.yml")
	assert.Error(t, err, "Should have error because of missing services")

	// reset
	configInfo = config{}
	err = makeConfig("config_test/config_no_recordpath.yml")
	assert.Error(t, err, "Should have error because of missing record_path")

	// reset
	configInfo = config{}
	err = makeConfig("config_test/config_no_hosts.yml")
	assert.Error(t, err, "Should have error because of missing es_hosts")

	// reset
	configInfo = config{}
	err = makeConfig("config_test/config_no_influxdb.yml")
	assert.NoError(t, err, "should not have error")
	assert.Equal(t, []string{"http://192.168.1.210:9200"}, configInfo.EsHosts, "es_hosts wrong")
	assert.Equal(t, []string{"a", "b"}, configInfo.Services, "services wrong")
	assert.Equal(t, "config_test/record_path.json", configInfo.Path, "record_path wrong")
	assert.Equal(t, "", configInfo.IFInfo.Host, "influxdb host wrong")

	// reset
	configInfo = config{}
	err = makeConfig("config_test/config.yml")
	assert.NoError(t, err, "should not have error")
	assert.Equal(t, []string{"http://192.168.1.210:9200"}, configInfo.EsHosts, "es_hosts wrong")
	assert.Equal(t, []string{"a", "b"}, configInfo.Services, "services wrong")
	assert.Equal(t, "config_test/record_path.json", configInfo.Path, "record_path wrong")
	assert.Equal(t, "http://192.168.1.1", configInfo.IFInfo.Host, "influxdb host wrong")
	assert.Equal(t, "abc", configInfo.IFInfo.Username, "influxdb username wrong")
	assert.Equal(t, "cba", configInfo.IFInfo.Password, "influxdb password wrong")
	assert.Equal(t, "test", configInfo.IFInfo.Database, "influxdb db wrong")
	// record info
	assert.Equal(t, "abc", records["a"], "record wrong")
}
