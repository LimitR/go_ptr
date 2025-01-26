# go_ptr

```shell
go get github.com/LimitR/go_ptr
```

```go
func main() {
    valuePtr := share_ptr.MakeShare("hello")

    fmt.Println(*valuePtr.Get()) // 'hello'
	
    valuePtr.Free()

    fmt.Println(valuePtr.Get()) // nil
}
```

```go
func main() {
    valuePtr := share_ptr.MakeShare("hello")

    fmt.Println(*valuePtr.Get()) // 'hello'

    someValuePtr := valuePtr.Copy()

    valuePtr.Free()

    *someValuePtr.Get() = "not hello"
    // Or
    someValuePtr.SetValue("not hello") // New pointer

    fmt.Println(*valuePtr.Get())     // 'not hello'
    fmt.Println(*someValuePtr.Get()) // 'not hello'

    someValuePtr.Free()
}
```
Use only `.Copy()`
```go
SomeFunc(myPtr.Copy())
```


## Array:
```go
func main() {
    array := share_ptr.MakeShareArray[int](0, 3)

    array.Append(1)
    array.Append(2)
    array.Append(3)

    *array.GetElement(2) = 4

    iter := array.Iter()

    for {
        v := iter()
        if v == nil {
            break
        }
        fmt.Println(*v) // 1, 2, 4
    }

    array.Free()

    fmt.Println(array.GetElement(2)) // nil
}
```

## GC
```go
valuePtr := share_ptr.MakeShare("hello").EnableGCHowFree()
valuePtr.Free() // Start goroutine with GC
```