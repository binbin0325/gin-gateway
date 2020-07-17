package endpoints

import (
	"testing"
	"unsafe"
)

func BenchmarkAssertDirectInvoke(b *testing.B) {
	s := new(Student)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DirectInvoke(s)
	}
	_ = s
}

func BenchmarkAssertPointerInvoke(b *testing.B) {
	s := new(Student)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PointerInvoke(unsafe.Pointer(s))
	}
	_ = s
}

func BenchmarkAssertInterfaceInvoke(b *testing.B) {
	s := new(Student)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InterfaceInvoke(s)
	}
	_ = s
}

func BenchmarkAssertInterfaceNotPointInvoke(b *testing.B) {
	s := Student{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InterfaceInvoke1(s)
	}
	_ = s
}
