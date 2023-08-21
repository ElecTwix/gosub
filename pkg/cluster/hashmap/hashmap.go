package hashmap

import "io"

type HashMapCluster struct {
	Subs map[string]io.WriteCloser
}

func NewHashMapCluster() *HashMapCluster {
	return &HashMapCluster{
		Subs: make(map[string]io.WriteCloser),
	}
}

func (a *HashMapCluster) AddSub(id string, sub io.WriteCloser) {
	a.Subs[id] = sub
}

func (a *HashMapCluster) Write(p []byte) (n int, err error) {
	for _, sub := range a.Subs {
		n, err = sub.Write(p)
		if err != nil {
			return
		}
	}
	return
}

func (a *HashMapCluster) WriteToSub(id string, p []byte) (n int, err error) {
	val, ok := a.Subs[id]
	if !ok {
		return 0, io.ErrClosedPipe
	}
	return val.Write(p)
}

func (a *HashMapCluster) RemoveSub(id string) {
	a.Subs[id].Close()
	delete(a.Subs, id)
}

func (a *HashMapCluster) Close() (err error) {
	for _, sub := range a.Subs {
		err = sub.Close()
		if err != nil {
			return
		}
	}
	return
}
