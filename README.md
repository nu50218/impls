# impls

[![Go Report Card](https://goreportcard.com/badge/github.com/nu50218/impls)](https://goreportcard.com/report/github.com/nu50218/impls)

(README書きかけ)

implsはinterfaceの実装を見つけたり、実装からinterfaceを見つけたりできるツールです。

## Install

`$ go get -u github.com/nu50218/impls`

## Usage

- interfaceから型を見つける

```go
$ impls types go/ast.Expr
/usr/local/go/src/go/ast/ast.go:411:2 ast.ArrayType
/usr/local/go/src/go/ast/ast.go:268:2 ast.BadExpr
/usr/local/go/src/go/ast/ast.go:288:2 ast.BasicLit
︙
```

- interfaceから変数を見つける

```go
$ impls vars error fmt
/usr/local/go/src/fmt/scan.go:466:5 fmt.boolError
/usr/local/go/src/fmt/scan.go:465:5 fmt.complexError
```

- 型からinterfaceを見つける

いずれのサブコマンドも第四引数以降にロードさせたいパッケージを渡すことができます。

```go
$ impls interfaces bytes.Buffer io
/usr/local/go/src/io/io.go:243:6 io.ByteReader
/usr/local/go/src/io/io.go:254:6 io.ByteScanner
/usr/local/go/src/io/io.go:260:6 io.ByteWriter
/usr/local/go/src/io/io.go:120:6 io.ReadWriter
/usr/local/go/src/io/io.go:77:6 io.Reader
/usr/local/go/src/io/io.go:170:6 io.ReaderFrom
/usr/local/go/src/io/io.go:269:6 io.RuneReader
/usr/local/go/src/io/io.go:280:6 io.RuneScanner
/usr/local/go/src/io/io.go:286:6 io.StringWriter
/usr/local/go/src/io/io.go:90:6 io.Writer
/usr/local/go/src/io/io.go:181:6 io.WriterTo
```
