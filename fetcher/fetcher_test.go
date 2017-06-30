package fetcher

import (
	"testing"
)

func TestGetMessage(t *testing.T) {
	// jCmd := fmt.Sprintf("journalctl -u dnsupdater -o json --after-cursor=\"%s\"", "s=942af3eb5dfe4814991a6d15f2ec61b7;i=5e7;b=f62dc60799c2434c8f996cd6ec22e24a;m=1292c8bad5;t=5411831a21d5e;x=7deba9a44c029328")
	// results, err := GetMessages("dnsupdater", "ssh", "leeq@192.168.1.201", jCmd)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// for _, one := range results {
	// 	t.Log("Message: ", string(one.Message))
	// }
}
