package ringmap

import "sync"

type RingMap struct {
	obhash      map[string][]byte
	fifo        []string
	fifoPointer int
	ringsize    int
	mutex       sync.Mutex
}

func NewRingMap(size int) *RingMap {
	r := &RingMap{}
	r.ringsize = size
	r.obhash = make(map[string][]byte)
	r.fifo = make([]string, size)
	r.fifoPointer = 0
	return r
}

func (r *RingMap) Get(key string) []byte {
	r.mutex.Lock()
	res := r.obhash[key]
	r.mutex.Unlock()
	return res
}

func (r *RingMap) Put(key string, value []byte) {
	r.mutex.Lock()
	r.obhash[key] = value
	r.fifo[r.fifoPointer] = key
	r.fifoPointer++
	if r.fifoPointer >= r.ringsize {
		r.fifoPointer = 0
	}
	deleteKey := r.fifo[r.fifoPointer]
	if deleteKey != "" {
		delete(r.obhash, deleteKey)
	}
	r.mutex.Unlock()
}
