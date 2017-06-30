package j2es

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordStatusWithoutInfluxdbSettings(t *testing.T) {
	configInfo = config{}
	err := makeConfig("config_test/config_no_influxdb.yml")
	assert.NoError(t, err, "should not have error")

	err = recordStatus(int64(0), false)
	assert.NoError(t, err, "should not have error")
}

func TestGetCommand(t *testing.T) {
	cmd, args := getCommand("abc", "cba")
	assert.Equal(t, "journalctl", cmd, "command is wrong")
	assert.EqualValues(t, []string{"-u", "abc", "-o", "json", "--after-cursor=cba"}, args, "arguments wrong")
}
