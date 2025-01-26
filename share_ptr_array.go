package share_ptr

import (
	"arena"
	"runtime"
	"unsafe"
)

type SharePtrArray[T any] struct {
	SharePtr[T]
	array               []T
	sizeElement         uintptr
	first               bool
	pointerFirstElement unsafe.Pointer
	iterMemory          *arena.Arena
}

func MakeShareArray[T any](len, cap int) *SharePtrArray[T] {
	memory := arena.NewArena()

	iterMemory := arena.NewArena()

	lenArray := 1
	if len > 1 {
		lenArray = len
	}

	if cap == 0 {
		cap = lenArray
	}

	array := arena.MakeSlice[T](memory, lenArray, cap)

	sizeElement := unsafe.Sizeof(array[0])

	pointer := unsafe.Pointer(&array)

	pointerFirstElement := unsafe.Pointer(&array[0])

	spa := &SharePtrArray[T]{
		SharePtr[T]{
			counter: 1,
			pointer: pointer,
			memory:  memory,
		},
		array,
		sizeElement,
		true,
		pointerFirstElement,
		iterMemory,
	}

	runtime.SetFinalizer(spa, func(a *SharePtrArray[T]) {
		a.memory.Free()
		a.iterMemory.Free()
	})

	return spa
}

func (s *SharePtrArray[T]) Append(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.first {
		s.array[0] = value
		s.first = false
		return
	}
	s.array = append(s.array, value)
}

// GetElement element by index
func (s *SharePtrArray[T]) GetElement(index int) *T {
	if s.pointer == nil || len(s.array)-1 < index {
		return nil
	}
	return (*T)(unsafe.Add(s.pointerFirstElement, uintptr(index)*s.sizeElement))
}

func (s *SharePtrArray[T]) Get() *[]T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return (*[]T)(s.pointer)
}

// SetValue not used
func (s *SharePtrArray[T]) SetValue(value T) {}

func (s *SharePtrArray[T]) Iter() func() *T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := arena.New[int](s.iterMemory)
	cb := func() *T {
		result := s.GetElement(*index)
		*index++
		return result
	}

	return cb
}

func (s *SharePtrArray[T]) Copy() *SharePtrArray[T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pointer == nil {
		return nil
	}
	s.counter++
	return s
}

func (s *SharePtrArray[T]) Free() {
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
			s.iterMemory.Free()
		}
		s.pointer = nil
		s.pointerFirstElement = nil
		s.first = false
		s.array = nil
	}
	s.counter--
}
