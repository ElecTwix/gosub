package gosub

import (
	"io"

	"github.com/ElecTwix/gosub/pkg/cluster"
)

type GoSub struct {
	ClusterMap map[string]cluster.Cluster
}

func NewGoSub() *GoSub {
	return &GoSub{
		ClusterMap: make(map[string]cluster.Cluster),
	}
}

func (g *GoSub) AddCluster(id string, c cluster.Cluster) {
	g.ClusterMap[id] = c
}

func (g *GoSub) RemoveCluster(id string) {
	g.ClusterMap[id].Close()
	delete(g.ClusterMap, id)
}

func (g *GoSub) AddSub(ClusterID, SubID string, writer io.WriteCloser) {
	g.ClusterMap[ClusterID].AddSub(SubID, writer)
}

func (g *GoSub) RemoveSub(ClusterID, SubID string) {
	g.ClusterMap[ClusterID].RemoveSub(SubID)
}

func (g *GoSub) WriteAll(p []byte) (n int, err error) {
	for _, c := range g.ClusterMap {
		n, err = c.Write(p)
		if err != nil {
			return
		}
	}
	return
}

func (g *GoSub) WriteToCluster(ClusterID string, p []byte) (n int, err error) {
	return g.ClusterMap[ClusterID].Write(p)
}

func (g *GoSub) WriteToSub(ClusterID, SubID string, p []byte) (n int, err error) {
	return g.ClusterMap[ClusterID].WriteToSub(SubID, p)
}

func (g *GoSub) Close() (err error) {
	for _, c := range g.ClusterMap {
		err = c.Close()
		if err != nil {
			return
		}
	}
	return
}

func (g *GoSub) CloseCluster(ClusterID string) (err error) {
	return g.ClusterMap[ClusterID].Close()
}
