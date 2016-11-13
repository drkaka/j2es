package uploader

type record struct {
	Timestamp string `json:"timestamp"`
	IsOK      bool   `json:"ok"`
}

// TestUploader to test uploader.
// func TestUploader(t *testing.T) {
// 	PrepareUploader([]string{"http://192.168.1.210:9200"})

// 	var one record
// 	one.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
// 	one.IsOK = true
// 	AddRecord("twitter", one)

// 	time.Sleep(10 * time.Millisecond)

// 	var two record
// 	two.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
// 	two.IsOK = false
// 	AddRecord("twitter", two)

// 	PostToES()
// }
