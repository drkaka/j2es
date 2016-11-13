package fetcher

import (
	"bufio"
	"encoding/json"
	"j2es/logger"
	"os/exec"
	"strconv"

	"github.com/uber-go/zap"
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
	Comm          string          `json:"_COMM"`
	Exe           json.RawMessage `json:"_EXE,omitempty"`
	CmdLine       json.RawMessage `json:"_CMDLINE,omitempty"`
	CapEffective  json.RawMessage `json:"_CAP_EFFECTIVE,omitempty"`
	SysdGroup     json.RawMessage `json:"_SYSTEMD_CGROUP,omitempty"`
	SysdUnit      json.RawMessage `json:"_SYSTEMD_UNIT,omitempty"`
	SysdSlice     json.RawMessage `json:"_SYSTEMD_SLICE,omitempty"`
	MachineID     json.RawMessage `json:"_MACHINE_ID,omitempty"`
	Hostname      json.RawMessage `json:"_HOSTNAME,omitempty"`
	Message       json.RawMessage `json:"MESSAGE"`
}

// GetMessages to get journal messages from a command.
func GetMessages(service, name string, arg ...string) ([]Result, error) {
	var results []Result
	cmd := exec.Command(name, arg...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return results, err
	}

	if err := cmd.Start(); err != nil {
		return results, err
	}
	logger.Log.Debug("CMD started.")

	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		var result JMessage

		b := scanner.Bytes()
		logger.Log.Debug("", zap.String("line", string(scanner.Bytes())))

		if err := json.Unmarshal(b, &result); err != nil {
			return results, err
		}

		// check whether the log come from applicaiton.
		if result.Comm != service {
			logger.Log.Debug("escape", zap.String("comm", result.Comm))
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

	// lines := bytes.Split(out, []byte("\n"))
	// logger.Log.Debug("", zap.Int("lines", len(lines)))

	// for i := 0; i < len(lines); i++ {
	// 	if len(lines[i]) == 0 {
	// 		continue
	// 	}

	// 	var result JMessage
	// 	if err := json.Unmarshal(lines[i], &result); err != nil {
	// 		return results, err
	// 	}

	// 	// check whether the log come from applicaiton.
	// 	if result.Comm != name {
	// 		logger.Log.Debug("escape", zap.String("comm", result.Comm))
	// 		continue
	// 	}

	// }
}
