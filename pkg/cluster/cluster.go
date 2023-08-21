package cluster

import "io"

type Cluster interface {
	AddSub(id string, sub io.Writer)
	RemoveSub(id string)
	Write(p []byte) (n int, err error)
	WriteToSub(id string, p []byte) (n int, err error)
	Close() (err error)
}
