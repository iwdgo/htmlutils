package parsing

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"log"
)

// ExampleIncludeNode is using the provided files to demonstrate use
func ExampleIncludedNode() {
	// f1 est une partie de f2
	toFind := html.Node{nil, nil, nil, nil, nil, html.ElementNode,
		0, "table", "",
		[]html.Attribute{{"", "class", "fixed"}},
	}
	pm, _ := ParseFile(f1)
	m := FindNode(pm, toFind) // searching <table> in d1
	if m == nil {
		fmt.Printf("%s not found in %s \n", PrintData(&toFind), f1)
	}

	pn, _ := ParseFile(f2)
	n := FindNode(pn, toFind) // searching <table> in d2
	if n == nil {
		fmt.Printf("%s not found in %s \n", PrintData(&toFind), f2)
	}
	// Is n included in m
	if f := IncludedNode(n, m); f != nil {
		fmt.Printf("nodes structures diverge from : %s\n", PrintData(f))
	}
	// Output:
}

// Examples below cannot be tested as multiple lines output fails (https://github.com/golang/go/issues/26460).
// They were prefixed by "Test" (including signature) to avoid failures. No test occurs.
// They demonstrate the output possibilities.

const HTMLf = `<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`

// Printing the node tree
func ExampleExploreNodeTags() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, err := html.Parse(b) // Only place where err of Parse is checked
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	ExploreNode(o, "", html.TextNode)
	// Output: HTML Fragment to compare against  (Text)
	//  others below (Text) to test  (Text)
	//  diffs (Text)
}
func ExampleExploreNodeAll() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	ExploreNode(o, "", html.ErrorNode)
	// Output: (Document)
	//  html (Element)
	//  head (Element) body (Element)
	//  p (Element) [{ class ex1}]
	//  HTML Fragment to compare against  (Text) em (Element)
	//  others below (Text) to test  (Text) sub (Element)
	//  diffs (Text)
}
func ExamplePrintTagswoSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	PrintTags(o, "", false) // +1,6%
	// Output:
	// (Document)
	//html (Element)
	//head (Element)
	//body (Element)
	//p (Element) [{ class ex1}]
	//HTML Fragment to compare against  (Text)
	//em (Element)
	//others below (Text)
	// to test  (Text)
	//sub (Element)
	//diffs (Text)
}

// Same as before but only tags stopping at a searched tag
func ExamplePrintTagswSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)    // err ignored as failure is detected before
	PrintTags(o, "em", true) //
	// Output:
	//html (Element)
	//head (Element)
	//body (Element)
	//p (Element) [{ class ex1}]
	//em (Element)
	//[em] found. Stopping exploration
	//sub (Element)
}

func ExamplePrintNodeswoSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	PrintNodes(o, nil, html.ErrorNode, 0)
	// Output: html (Element)
	//. head (Element) body (Element)
	//.. p (Element) [{ class ex1}]
	//... HTML Fragment to compare against  (Text) em (Element)
	//.... others below (Text) to test  (Text) sub (Element)
	//.... diffs (Text)
}

func ExamplePrintNodeswSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)

	var tagToFind html.Node
	tagToFind.Type = html.ElementNode
	tagToFind.Data = "p"
	tagToFind.Attr = []html.Attribute{{"", "class", "ex1"}}

	PrintNodes(o, &tagToFind, html.ErrorNode, 0)
	// Output: html (Element)
	//. head (Element) body (Element)
	//.. p (Element) [{ class ex1}]
	//tag found: p (Element) [{ class ex1}]
	//... HTML Fragment to compare against  (Text) em (Element)
	//.... others below (Text) to test  (Text) sub (Element)
	//.... diffs (Text)
}
