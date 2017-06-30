package fetcher

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/drkaka/lg"
	"go.uber.org/zap"
)

// Result to define a journal message result.
type Result struct {
	Cursor  string
	Message json.RawMessage
}

// JMessage to define journal a message.
type JMessage struct {
	Cursor        string          `json:"__CURSOR"`
	RealTime      json.RawMessage `json:"__REALTIME_TIMESTAMP,omitempty"`
	MonotonicTime json.RawMessage `json:"__MONOTONIC_TIMESTAMP,omitempty"`
	BootID        json.RawMessage `json:"_BOOT_ID,omitempty"`
	Transport     json.RawMessage `json:"_TRANSPORT,omitempty"`
	Priority      json.RawMessage `json:"PRIORITY,omitempty"`
	SyslFacility  json.RawMessage `json:"SYSLOG_FACILITY,omitempty"`
	SyslID        json.RawMessage `json:"SYSLOG_IDENTIFIER,omitempty"`
	PID           json.RawMessage `json:"_PID,omitempty"`
	UID           json.RawMessage `json:"_UID,omitempty"`
	GID           json.RawMessage `json:"_GID,omitempty"`
	Comm          json.RawMessage `json:"_COMM,omitempty"`
	Exe           json.RawMessage `json:"_EXE,omitempty"`
	CmdLine       json.RawMessage `json:"_CMDLINE,omitempty"`
	CapEffective  json.RawMessage `json:"_CAP_EFFECTIVE,omitempty"`
	SysdGroup     json.RawMessage `json:"_SYSTEMD_CGROUP,omitempty"`
	SysdUnit      string          `json:"_SYSTEMD_UNIT"`
	SysdSlice     json.RawMessage `json:"_SYSTEMD_SLICE,omitempty"`
	MachineID     json.RawMessage `json:"_MACHINE_ID,omitempty"`
	Hostname      json.RawMessage `json:"_HOSTNAME,omitempty"`
	Message       json.RawMessage `json:"MESSAGE"`
}

// GetMessages to get journal messages from a command.
func GetMessages(service, cmdName string, arg ...string) ([]Result, error) {
	var results []Result
	cmd := exec.Command(cmdName, arg...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return results, err
	}

	if err := cmd.Start(); err != nil {
		return results, err
	}
	lg.L(nil).Debug("CMD started.")

	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		var result JMessage

		b := scanner.Bytes()
		lg.L(nil).Debug("", zap.String("line", string(scanner.Bytes())))

		if err := json.Unmarshal(b, &result); err != nil {
			return results, err
		}

		// check whether the log comes from the service.
		// this will omit system messages like start or stop
		if result.SysdUnit != strings.Join([]string{service, "service"}, ".") {
			lg.L(nil).Debug("escape", zap.String("unit", result.SysdUnit))
			continue
		}

		var one Result
		one.Cursor = result.Cursor

		real, err := strconv.Unquote(string(result.Message))
		if err != nil {
			return results, err
		}
		one.Message = []byte(real)

		results = append(results, one)
	}
	return results, nil
}
