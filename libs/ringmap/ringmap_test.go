package ringmap

import "testing"

func Test_Ring(t *testing.T) {
	r := NewRingMap(5)

	r.Put("a", []byte("1"))
	r.Put("b", []byte("2"))
	r.Put("c", []byte("3"))
	r.Put("d", []byte("4"))
	r.Put("e", []byte("5"))
	r.Put("f", []byte("6"))
	r.Put("g", []byte("7"))
	r.Put("h", []byte("8"))
	r.Put("i", []byte("9"))
	r.Put("j", []byte("10"))

	print(string(r.Get("d")))
}
