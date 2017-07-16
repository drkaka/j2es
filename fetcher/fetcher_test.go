package fetcher

import (
	"fmt"
	"testing"

	"github.com/drkaka/lg"
)

func TestGetMessage(t *testing.T) {
	lg.InitLogger(true)

	jCmd := fmt.Sprintf("journalctl -u longlog -o json")
	results, err := GetMessages("longlog", "ssh", "leeq@192.168.1.201", jCmd)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("got messages:", len(results))
	for _, one := range results {
		fmt.Println("Message: ", string(one.Message))
	}
}
