package gosub_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/ElecTwix/gosub"
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

func NewGoSub(t *testing.T) {
	t.Parallel()

	t.Run("NewGoSub", func(t *testing.T) {
		gosub := gosub.NewGoSub()

		if gosub == nil {
			t.Fatal("gosub is nil")
		}

		if gosub.ClusterMap == nil {
			t.Error("gosub.ClusterMap is nil")
		}

		if len(gosub.ClusterMap) != 0 {
			t.Error("gosub.ClusterMap is not empty")
		}

	})
}

func TestGoSubWithHashMapCluster_All(t *testing.T) {
	t.Parallel()

	t.Run("TestHashMapCluster_All", func(t *testing.T) {

		gosub := gosub.NewGoSub()
		cluster := hashmap.NewHashMapCluster()
		gosub.AddCluster("test", cluster)

		writerArr := make([]io.WriteCloser, 0)
		for i := 0; i < 10; i++ {
			writerArr = append(writerArr, &testWriter{})
		}

		for i, writer := range writerArr {
			gosub.AddSub("test", fmt.Sprintf("test%d", i), writer)
		}

		if gosub == nil {
			t.Fatal("gosub is nil")
		}

		if gosub.ClusterMap == nil {
			t.Error("gosub.ClusterMap is nil")
		}

		if len(gosub.ClusterMap) != 1 {
			t.Error("gosub.ClusterMap is not empty")
		}

		if len(cluster.Subs) != 10 {
			t.Error("cluster.Subs is not empty")
		}

		if len(writerArr) != 10 {
			t.Error("writerArr is not empty")
		}

		for i, writer := range writerArr {
			if len(writer.(*testWriter).data) != 0 {
				t.Errorf("writerArr[%d].data is not empty", i)
			}
		}

		for i, writer := range writerArr {
			n, err := writer.Write([]byte("test"))
			if err != nil {
				t.Errorf("writerArr[%d].Write returned error: %v", i, err)
			}
			if n != 4 {
				t.Errorf("writerArr[%d].Write returned n: %d", i, n)
			}
		}

		for i, writer := range writerArr {
			gosub.RemoveSub("test", fmt.Sprintf("test%d", i))
			if !writer.(*testWriter).closed {
				t.Errorf("writerArr[%d] is not closed", i)
			}
		}

		if len(cluster.Subs) != 0 {
			t.Error("cluster.Subs is not empty")
		}

		gosub.RemoveCluster("test")

		if len(gosub.ClusterMap) != 0 {
			t.Error("gosub.ClusterMap is not empty")
		}
	})

}
