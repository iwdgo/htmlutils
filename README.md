[![GoDoc](https://godoc.org/github.com/iWdGo/htmlutils?status.svg)](https://godoc.org/github.com/iWdGo/htmlutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwdgo/htmlutils)](https://goreportcard.com/report/github.com/iwdgo/htmlutils)
[![codecov](https://codecov.io/gh/iWdGo/htmlutils/branch/master/graph/badge.svg)](https://codecov.io/gh/iWdGo/htmlutils)

# Exploring html.Node trees

html.Node trees as parsed by [golang.org/x/net/html](https://godoc.org/golang.org/x/net/html).
Basic search and comparison of tags or nodes is restricted by HTML rules and parsing behaviour.
The search of an HTML tag using a *node.HTML type can be done using any non-pointer value.
The first match is always returned.

HTML.node trees and sub-trees can be compared.
Text value of a tag like a title or an error message can be checked.

# Good to know

Siblings must always have the same order or comparison fails.
Order of attributes is irrelevant.
