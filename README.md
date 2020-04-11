# uber-go/zap go-pkgz/lgr bridge

## goal

the main goal of this package is used by [ttys3/remark42](https://github.com/ttys3/remark42)
to replace the used
`github.com/go-pkgz/lgr` logger with [go.uber.org/zap](https://github.com/uber-go/zap)

## usage

add below config in your `go.mod`

```bash
replace (
    github.com/go-pkgz/lgr v0.6.3 => github.com/ttys3/lgr master
)
```

then run:

```bash
go mod tidy
```
