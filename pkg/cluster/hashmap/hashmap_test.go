package hashmap_test

import (
	"io"
	"testing"

	"github.com/ElecTwix/gosub/pkg/cluster/hashmap"
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

func TestNewHashMapCluster(t *testing.T) {
	t.Parallel()

	t.Run("TestNewHashMapCluster", func(t *testing.T) {
		cluster := hashmap.NewHashMapCluster()
		if cluster == nil {
			t.Fatal("cluster is nil")
		}

		if cluster.Subs == nil {
			t.Error("cluster.Subs is nil")
		}

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}
	})
}

func TestHashMapCluster_AddSub(t *testing.T) {
	t.Parallel()

	t.Run("TestHashMapCluster_All", func(t *testing.T) {
		cluster := hashmap.NewHashMapCluster()
		if cluster == nil {
			t.Fatal("cluster is nil")
		}

		if cluster.Subs == nil {
			t.Error("cluster.Subs is nil")
		}

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}

		TestWriter := &testWriter{}

		cluster.AddSub("test", TestWriter)

		if len(cluster.Subs) != 1 {
			t.Error("cluster.Subs is not empty")
		}

		if len(TestWriter.data) != 0 {
			t.Error("TestWriter.data is not empty")
		}

		_, err := cluster.Write([]byte("test"))
		if err != nil {
			t.Error(err)
		}

		if len(TestWriter.data) != 4 {
			t.Error("TestWriter.data is not empty")
		}

		if string(TestWriter.data) != "test" {
			t.Error("TestWriter.data is not equal to test")
		}

		_, err = cluster.WriteToSub("test", []byte("test"))
		if err != nil {
			t.Error(err)
		}

		if len(TestWriter.data) != 8 {
			t.Error("TestWriter.data is not empty")
		}

		if string(TestWriter.data) != "testtest" {
			t.Error("TestWriter.data is not equal to testtest")
		}

		cluster.RemoveSub("test")

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}

		if !TestWriter.closed {
			t.Error("TestWriter is not closed")
		}

		_, err = cluster.WriteToSub("test", []byte("test"))
		if err == nil {
			t.Error("cluster.WriteToSub should return error")
		}

		cluster.Close()

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}

		if !TestWriter.closed {
			t.Error("TestWriter is not closed")
		}

	})

}
