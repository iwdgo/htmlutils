package parsing

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"testing"
)

// ExampleIncludeNode is using the provided files to demonstrate use
func ExampleIncludeNode() {
	// f1 est une partie de f2
	toFind := html.Node{nil, nil, nil, nil, nil, html.ElementNode,
		0, "table", "",
		[]html.Attribute{{"", "class", "fixed"}},
	}
	m := FindNode(ParseFile(f1), toFind) // searching <table> in d1
	if m == nil {
		fmt.Printf("%s not found in %s \n", printData(&toFind), f1)
	}

	n := FindNode(ParseFile(f2), toFind) // searching <table> in d2
	if n == nil {
		fmt.Printf("%s not found in %s \n", printData(&toFind), f2)
	}
	// Is n included in m
	if f := IncludedNode(n, m); f != nil {
		fmt.Printf("nodes structures diverge from : %s\n", printData(f))
	}
	// Output:
}

// Examples below cannot be tested as multiple lines output fails (https://github.com/golang/go/issues/26460).
// They were prefixed by "Test" (including signature) to avoid failures. No test occurs.
// They demonstrate the output possibilities.

const HTMLf = `<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`

// Printing the node tree
func TestExampleExploreNodeTags(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, err := html.Parse(b)
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	ExploreNode(o, "", html.TextNode)
	// Output: HTML Fragment to compare against  (Text)
	// others below (Text)       to test  (Text)
	// diffs (Text)
}
func TestExampleExploreNodeAll(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, err := html.Parse(b)
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	ExploreNode(o, "", html.ErrorNode)
	// Output:
	//  (Document)
	// html (Element)
	// head (Element)          body (Element)
	// p (Element) [{ class ex1}]
	// HTML Fragment to compare against  (Text)        em (Element)
	// others below (Text)      to test  (Text)        sub (Element)
	// diffs (Text)
}
func TestExamplePrintTagswoSearch(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, err := html.Parse(b)
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
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

//
func TestExamplePrintTagswSearch(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, `<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`)
	o, err := html.Parse(b)
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	PrintTags(o, "em", false) // +1,6%
	// Output:
	// (Document)
	//html (Element)
	//head (Element)
	//body (Element)
	//p (Element) [{ class ex1}]
	//HTML Fragment to compare against  (Text)
	//em (Element)
	//[em] found. Stopping exploration
	// to test  (Text)
	//sub (Element)
	//diffs (Text)
}

func TestExamplePrintNodes(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, `<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`)
	o, err := html.Parse(b)
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	PrintNodes(o, nil, html.ErrorNode, 0)
	// Output: html (Element)
	//.head (Element) body (Element)
	//..p (Element) [{ class ex1}]
	//...HTML Fragment to compare against  (Text) em (Element)
	//....others below (Text)  to test  (Text) sub (Element)
	//....diffs (Text)
}
