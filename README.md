# confmaker

Generate a configuration source file based on a `go` struct

```bash
go build
./confmaker -in=foo.go
```

foo.go

```go
package foo

type Foo struct {
    Field1 string
    Field2 int64
}
```