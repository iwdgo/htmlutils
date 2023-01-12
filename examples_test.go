package parsing_test

import (
	"bytes"
	"fmt"
	parsing "github.com/iwdgo/htmlutils"
	"golang.org/x/net/html"
	"log"
)

const HTMLf = `<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`

func ExampleGetText() {
	b := new(bytes.Buffer)
	_, _ = fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b) // Any parsing error would occured elsewhere
	w := new(bytes.Buffer)
	parsing.GetText(o, w)
	if s := fmt.Sprint(w); s != "HTML Fragment to compare against others below to test diffs" {
		fmt.Println("incorrect text")
	}
}

// ExampleExploreNode_tags only prints text.
func ExampleExploreNode_tags() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, err := html.Parse(b) // Only place where err of Parse is checked
	if err != nil {
		log.Fatalf("parsing error:%v\n", err)
	}
	parsing.ExploreNode(o, "", html.TextNode)
	// Output: HTML Fragment to compare against  (Text)
	//  others below (Text) to test  (Text)
	//  diffs (Text)
}

// ExampleExploreNode_all prints the complete node tree.
func ExampleExploreNode_all() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	parsing.ExploreNode(o, "", html.ErrorNode)
	// Output: (Document)
	//  html (Element)
	//  head (Element) body (Element)
	//  p (Element) [{ class ex1}]
	//  HTML Fragment to compare against  (Text) em (Element)
	//  others below (Text) to test  (Text) sub (Element)
	//  diffs (Text)
}

// ExamplePrintTags_woSearch is not using the search part.
func ExamplePrintTags_woSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	parsing.PrintTags(o, "", false) // +1,6%
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

// ExamplePrintTags_wSearch is the previous example stopping at a searched tag
func ExamplePrintTags_wSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)            // err ignored as failure is detected before
	parsing.PrintTags(o, "em", true) //
	// Output:
	//html (Element)
	//head (Element)
	//body (Element)
	//p (Element) [{ class ex1}]
	//em (Element)
	//[em] found. Stopping exploration
	//sub (Element)
}

// ExamplePrintNodes_woSearch prints all nodes without using search.
func ExamplePrintNodes_woSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)
	parsing.PrintNodes(o, nil, html.ErrorNode, 0)
	// Output: html (Element)
	//. head (Element) body (Element)
	//.. p (Element) [{ class ex1}]
	//... HTML Fragment to compare against  (Text) em (Element)
	//.... others below (Text) to test  (Text) sub (Element)
	//.... diffs (Text)
}

// ExamplePrintNodes_wSearch is the previous example stopping at a searched node.
func ExamplePrintNodes_wSearch() {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b)

	var tagToFind html.Node
	tagToFind.Type = html.ElementNode
	tagToFind.Data = "p"
	tagToFind.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex1"}}

	parsing.PrintNodes(o, &tagToFind, html.ErrorNode, 0)
	// Output: html (Element)
	//. head (Element) body (Element)
	//.. p (Element) [{ class ex1}]
	//tag found: p (Element) [{ class ex1}]
	//... HTML Fragment to compare against  (Text) em (Element)
	//.... others below (Text) to test  (Text) sub (Element)
	//.... diffs (Text)
}
