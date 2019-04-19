[![GoDoc](https://godoc.org/github.com/iWdGo/htmlutils?status.svg)](https://godoc.org/github.com/iWdGo/htmlutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwdgo/htmlutils)](https://goreportcard.com/report/github.com/iwdgo/htmlutils)

# html.Node trees

This module provides basic tools to search and compare html.Node trees as produces by [golang.org/x/net/html](https://godoc.org/golang.org/x/net/html).
The original HTML page can be searched after parsing.
The search of a tag can be done using its name or any non-pointer attribute.
Trees and sub-trees can be compared.

# Good to know

Siblings must always have the same order or comparison fails.
Order of attributes is irrelevant.

# Testing

Test coverage is 98.3% which is the maximum as testing file existence is unfeasible w/o panic.