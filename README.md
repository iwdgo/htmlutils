[![Go Reference](https://pkg.go.dev/badge/github.com/iwdgo/htmlutils.svg)](https://pkg.go.dev/github.com/iwdgo/htmlutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwdgo/htmlutils)](https://goreportcard.com/report/github.com/iwdgo/htmlutils)
[![codecov](https://codecov.io/gh/iWdGo/htmlutils/branch/master/graph/badge.svg)](https://codecov.io/gh/iWdGo/htmlutils)

[![Build Status](https://app.travis-ci.com/iwdgo/htmlutils.svg?branch=master)](https://app.travis-ci.com/iwdgo/htmlutils)
[![Build Status](https://api.cirrus-ci.com/github/iwdgo/htmlutils.svg)](https://cirrus-ci.com/github/iwdgo/htmlutils)
[![Build status](https://ci.appveyor.com/api/projects/status/v6ce70t0jmqgehpw?svg=true)](https://ci.appveyor.com/project/iWdGo/htmlutils)
![Build status](https://github.com/iwdgo/htmlutils/workflows/Go/badge.svg)

# Exploring HTML structure

HTML is parsed using [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) which produces a tree.

The module provides basic functionality to compare HTML tags or nodes and their trees.
The search of an HTML tag using a `*node.HTML` type ignores pointers.
It always returns the first match. By ignoring some properties, tags like `<button>` are easy to count.
Text value of a tag (title, error message,...) can be checked.

# Good to know

Parsing is not done according to the complete syntax checker of HTML.
For instance, tags like `<p>` for which a closing tag would fail a comparison.

Siblings must always have the same order or comparison fails.
Order of attributes is treated as irrelevant.

# How to start

Detailed [documentation](https://pkg.go.dev/github.com/iwdgo/htmlutils) includes examples.

# Versions

`v1.0.6` updates golang/go/x/net package to remove CVE-2022-27664 which does not affect x/net/html 
`v1.0.5` requires Go 1.16+ as ioutil package use is removed.  
`v1.0.4` requires Go 1.17+ which implements lazy loading of modules to avoid go.mod updates.  
`v1.0.0` was created on Go 1.12 which supports modules.

