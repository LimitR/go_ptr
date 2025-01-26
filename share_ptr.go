package share_ptr

import (
	"arena"
	"runtime"
	"sync"
	"unsafe"
)

type SharePtr[T any] struct {
	mutex   sync.RWMutex
	memory  *arena.Arena
	pointer unsafe.Pointer
	counter int
	gc      bool
}

func MakeShare[T any](value T) *SharePtr[T] {
	memory := arena.NewArena()
	v := arena.New[T](memory)
	*v = value
	pointer := unsafe.Pointer(v)

	sp := &SharePtr[T]{
		counter: 1,
		pointer: pointer,
		memory:  memory,
	}

	runtime.SetFinalizer(sp, func(p *SharePtr[T]) {
		if p.gc {
			p.memory.Free()
		}
	})

	return sp
}

func (s *SharePtr[T]) Free() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pointer == nil {
		return
	}
	if s.counter-1 <= 0 {
		if s.gc {
			go runtime.GC()
		} else {
			s.memory.Free()
		}
		s.pointer = nil
	}
	s.counter--
}

func (s *SharePtr[T]) EnableGCHowFree() *SharePtr[T] {
	s.gc = true
	return s
}

func (s *SharePtr[T]) Copy() *SharePtr[T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pointer == nil {
		return nil
	}
	s.counter++
	return s
}

func (s *SharePtr[T]) Get() *T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pointer == nil {
		return nil
	}

	return (*T)(s.pointer)
}

func (s *SharePtr[T]) SetValue(value T) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.pointer == nil {
		return
	}

	v := arena.New[T](s.memory)
	*v = value
	s.pointer = unsafe.Pointer(v)
}
