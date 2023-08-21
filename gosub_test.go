package gosub_test

import (
	"fmt"
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

func TestHashMapCluster_All(t *testing.T) {
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

		writerArr := make([]io.WriteCloser, 0)
		for i := 0; i < 10; i++ {
			writerArr = append(writerArr, &testWriter{})
		}

		for i := 0; i < 10; i++ {
			cluster.AddSub(fmt.Sprintf("%d", i), writerArr[i])
		}

		if len(cluster.Subs) != 10 {
			t.Error("cluster.Subs is not empty")
		}

		for i := 0; i < 10; i++ {
			if len(writerArr[i].(*testWriter).data) != 0 {
				t.Error("TestWriter.data is not empty")
			}
		}

		for i := 0; i < 10; i++ {
			_, err := cluster.Write([]byte("test"))
			if err != nil {
				t.Error(err)
			}
		}

		for i := 0; i < 10; i++ {
			if len(writerArr[i].(*testWriter).data) != 40 {
				t.Error("TestWriter.data is not empty")
			}
		}

		for i := 0; i < 10; i++ {
			cluster.RemoveSub(fmt.Sprintf("%d", i))
			n, err := cluster.WriteToSub(fmt.Sprintf("%d", i), []byte("test"))
			if err == nil {
				t.Error("cluster.WriteToSub should return error")
			}
			if n != 0 {
				t.Error("cluster.WriteToSub should return 0")
			}
			if len(writerArr[i].(*testWriter).data) != 40 {
				t.Error("TestWriter.data is not empty")
			}
		}

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}

		for i := 0; i < 10; i++ {
			if !writerArr[i].(*testWriter).closed {
				t.Error("TestWriter is not closed")
			}
		}

	})

}
