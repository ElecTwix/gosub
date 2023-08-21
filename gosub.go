package gosub

import (
	"io"

	"github.com/ElecTwix/gosub/pkg/cluster"
)

type GoSub struct {
	clusterMap map[string]cluster.Cluster
}

func NewGoSub() *GoSub {
	return &GoSub{
		clusterMap: make(map[string]cluster.Cluster),
	}
}

func (g *GoSub) AddCluster(id string, c cluster.Cluster) {
	g.clusterMap[id] = c
}

func (g *GoSub) RemoveCluster(id string) {
	g.clusterMap[id].Close()
	delete(g.clusterMap, id)
}

func (g *GoSub) AddSub(ClusterID, SubID string, writer io.WriteCloser) {
	g.clusterMap[ClusterID].AddSub(SubID, writer)
}

func (g *GoSub) RemoveSub(ClusterID, SubID string) {
	g.clusterMap[ClusterID].RemoveSub(SubID)
}

func (g *GoSub) WriteAll(p []byte) (n int, err error) {
	for _, c := range g.clusterMap {
		n, err = c.Write(p)
		if err != nil {
			return
		}
	}
	return
}

func (g *GoSub) WriteToCluster(ClusterID string, p []byte) (n int, err error) {
	return g.clusterMap[ClusterID].Write(p)
}

func (g *GoSub) WriteToSub(ClusterID, SubID string, p []byte) (n int, err error) {
	return g.clusterMap[ClusterID].WriteToSub(SubID, p)
}

func (g *GoSub) Close() (err error) {
	for _, c := range g.clusterMap {
		err = c.Close()
		if err != nil {
			return
		}
	}
	return
}
