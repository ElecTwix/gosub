// GoSub is publish and subscribe system for Go.

// ## Installation
```bash
go get github.com/ElecTwix/gosub
```

// ## Usage
```go
package main

import (
  "fmt"
  "github.com/ElecTwix/gosub"
)

type testWriter struct {
	data   []byte
	closed bool
}

func (t *testWriter) Write(p []byte) (n int, err error) {
	if t.closed {
		return 0, io.ErrClosedPipe
	}
	t.data = append(t.data, p...)
	return len(p), nil
}

func (t *testWriter) Close() (err error) {
	t.closed = true
	return
}

func main() {
  // Create new instance of GoSub
	gosub := gosub.NewGoSub()
	cluster := hashmap.NewHashMapCluster()
	gosub.AddCluster("test_cluster", cluster)


  // Create new writer
  writer := &testWriter{}

  // Subscribe writer to topic to cluster
  gosub.Subscribe("test_Sub", writer)

  // Send message to all cluster subscribers
  gosub.WriteAll([]byte("Hello world!"))

  // Send message to all subscribers that in cluster
  gosub.Write("test_Sub", []byte("Hello world!"))

  // Send Only one sub message that in cluster
  gosub.WriteToSub("test_Sub", []byte("Hello world!"))

  // Unsubscribe writer from topic that in cluster
  gosub.RemoveSub("test_Sub", writer)

  // Close all subscribers that in cluster
  gosub.Close()

  // Close all subscribers that in cluster and remove cluster
  gosub.CloseCluster("test_cluster")

}
```
