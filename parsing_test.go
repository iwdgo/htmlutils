package parsing

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"testing"
)

const (
	f1 = "want_Test1.html"
	f2 = "want_Test2.html"
)

func TestAreEqual(t *testing.T) {
	var tag1, tag2 html.Node
	if !Equal(nil, nil) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	// nil == nil
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	tag1.Data = "h1"
	tag2.Data = "h2"
	// .Data != .Data
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	tag1.Type = html.ElementNode
	tag2.Data = tag1.Data
	// .Type != .Type
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	tag2.Type = tag1.Type
	tag2.Namespace = "ns"
	// .Namespace != .Namespace
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	tag1.Namespace = tag2.Namespace
	tag2.Attr = []html.Attribute{{"", "class", "fixed"}}
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	tag1.Attr = tag2.Attr
	// .Attr[0] == .Attr[0]
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	// .Attr[1] != .Attr[1]
	tag2.Attr = append(tag1.Attr, html.Attribute{"", "style", "h2"})
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	// .Attr == .Attr
	tag1.Attr = tag2.Attr
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	tag1.Attr = []html.Attribute{tag1.Attr[1], tag1.Attr[0]}
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
}

func TestPrintData(t *testing.T) {
	var tagToFind html.Node
	tagToFind.Type = html.ElementNode
	tagToFind.Data = "p"
	tagToFind.Attr = []html.Attribute{{"", "class", "ex2"}}
	want := "p (Element) [{ class ex2}]"
	if s := PrintData(&tagToFind); s != want {
		t.Errorf("printData: got %s, want %s", s, want)
	}
}

func TestGetText(t *testing.T) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, HTMLf)
	o, _ := html.Parse(b) // Any parsing error would occurred elsewhere
	w := new(bytes.Buffer)
	GetText(o, w)
	if s := fmt.Sprint(w); s != "HTML Fragment to compare against others below to test diffs" {
		t.Errorf("incorrect text")
	}
}

// Testing the search of a tag
func TestFindTag(t *testing.T) {
	d := ParseFile(f2)
	tagToFind := "table"
	want := tagToFind + " (Element) [{ class fixed}]" // PrintData value
	if n := FindTag(d, tagToFind, html.ElementNode); n == nil {
		t.Errorf("<%s> not found in %s\n", tagToFind, f2)
	} else if got := PrintData(n); got != want {
		t.Errorf("<%s> tag found is different: got %s, want %s\n", tagToFind, got, want)
	}
	tagToFind = "display"
	if n := FindTag(d, tagToFind, html.ElementNode); n != nil {
		t.Errorf("<%s> found in %s\n", tagToFind, f2)
	}
}

// Testing the the search of a node based on values which are not pointers, i.e. the tree of the node is not taken into
// account
func TestFindNodes(t *testing.T) {
	f := f2
	tagsToFind := []struct {
		n html.Node
		f bool
	}{
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "table", "",
			[]html.Attribute{{"", "class", "fixed"}}}, true},
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "p", "",
			[]html.Attribute{{"", "class", "ex2"}}}, true},
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "p", "",
			[]html.Attribute{{"", "class", "not-found"}}}, false},
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "p", "ns", // +0.9%
			[]html.Attribute{{"", "class", "not-found"}}}, false},
	}

	for _, m := range tagsToFind {
		want := PrintData(&m.n)
		if o := FindNode(ParseFile(f), m.n); o == nil && m.f { // Should be found
			t.Errorf("<%s> not found in %s\n", PrintData(&m.n), f)
		} else if o == nil && !m.f {
			// Not found is expected
		} else if got := PrintData(o); got != want {
			t.Errorf("tag found has differences: got %s, want %s\n", got, want)
		}
	}
}

// Testing nodes comparison
// Using nil to include (+0.8%)
func TestIncludeNode(t *testing.T) {
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "p"
	n.Attr = []html.Attribute{{"", "class", "ex2"}}
	if !Equal(IncludedNode(nil, &n), &n) {
		t.Errorf("nil does not include any node")
	}
}

// Using nil to include (+0.9%)
func TestIncludedNodeTyped(t *testing.T) {
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "p"
	n.Attr = []html.Attribute{{"", "class", "ex2"}}
	if !Equal(IncludedNodeTyped(nil, &n, html.ErrorNode), &n) {
		t.Errorf("nil does not include any node")
	}
}

// Using nil to include (+0.8%)
func TestIdenticalNilNodes(t *testing.T) {
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "p"
	n.Attr = []html.Attribute{{"", "class", "ex2"}}
	if !Equal(IdenticalNodes(nil, &n, html.ErrorNode), &n) {
		t.Errorf("nil does not include any node")
	}
}

func TestIncludeNodes(t *testing.T) {
	fragments := []struct {
		s     string
		equal bool
	}{
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`,
			true, // identical to itself
		},
		{`<p class="ex1">HTML Fragment to compare against <em>other below</em> to test <sub>diffs</sub></p>`,
			false, // identical to itself
		},
		// missing sub-node <sub>
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`, false},
		// additionnal sibling node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p><p class="ex2">outside</p>`,
			true},
		// additionnal sub-node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub><p class="ex2">inside</p></p>`,
			true},
	}
	b := new(bytes.Buffer)
	fmt.Fprint(b, fragments[0].s)
	original, err := html.Parse(b)
	// Buffer is empty after parsing
	if err != nil {
		t.Errorf("%v while parsing %s", err, fragments[0].s)
	} else {
		// No Sibling in the document after parsing
		// fmt.Printf("F %v\nL %v\nN %v\nP %v\n", original.FirstChild, original.LastChild, original.NextSibling, original.PrevSibling)
		ExploreNode(original, "", html.ErrorNode)
		fmt.Println("---")
	}
	for i, f := range fragments {
		fmt.Fprint(b, f.s)
		n, err := html.Parse(b)
		if err != nil {
			t.Errorf("---FAIL %d: %v while parsing %s", i, err, f.s)
		}
		if r := IncludedNode(original, n); r != nil && f.equal {
			// if r := compareNodeRecursive(original, n); r != nil {
			t.Errorf("---FAIL(%d): %s differs from [%s]", i, f.s, PrintData(r))
			// exploreNode(r, "", html.ErrorNode)
		} else if r == nil && !f.equal {
			t.Errorf("---FAIL(%d): no difference found with %s", i, f.s)
		} else if r != nil && !f.equal {
			fmt.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
		}
	}
}

// Testing tags of one type
func TestIncludeNodeTyped(t *testing.T) {
	fragments := []struct {
		s     string
		t     html.NodeType
		equal bool
	}{
		// identical to itself
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`,
			html.ErrorNode, // identical on a string basis
			true,
		},
		// text change
		{`<p class="ex1">HTML Fragment to compare against <em>other below</em> to test <sub>diffs</sub></p>`,
			html.ElementNode,
			true,
		},
		// missing sub-node <sub>
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`,
			html.ElementNode, false},
		// missing sub-node <sub> and trees do not have the same structure
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`,
			html.TextNode, false},
		// additionnal sibling node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p><p class="ex2">outside</p>`,
			html.ElementNode,
			true},
		// additionnal sub-node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub><p class="ex2">inside</p></p>`,
			html.ElementNode,
			true},
	}
	b := new(bytes.Buffer)
	fmt.Fprint(b, fragments[0].s)
	original, err := html.Parse(b)
	// Buffer is empty after parsing
	if err != nil {
		t.Errorf("%v while parsing %s", err, fragments[0].s)
	} else {
		// No Sibling in the document after parsing
		ExploreNode(original, "", original.Type)
		fmt.Println("---")
	}
	for i, f := range fragments {
		fmt.Fprint(b, f.s)
		n, err := html.Parse(b)
		if err != nil {
			t.Errorf("---FAIL %d: %v while parsing %s", i, err, f.s)
		}
		if r := IncludedNodeTyped(original, n, f.t); r != nil && f.equal {
			// if r := compareNodeRecursive(original, n); r != nil {
			t.Errorf("---FAIL(%d): %s differs from [%s]", i, f.s, PrintData(r))
			// exploreNode(r, "", html.ErrorNode)
		} else if r == nil && !f.equal {
			t.Errorf("---FAIL(%d): no difference found with %s", i, f.s)
		} else if r != nil && !f.equal {
			fmt.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
		}
	}
}

// Testing tags of one type
func TestIdenticalNodes(t *testing.T) {
	fragments := []struct {
		s     string
		t     html.NodeType
		equal bool
	}{
		// identical to itself
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`,
			html.ErrorNode, // identical on a string basis
			true,
		},
		// text change
		{`<p class="ex1">HTML Fragment to compare against <em>other below</em> to test <sub>diffs</sub></p>`,
			html.ElementNode,
			true,
		},
		// missing sub-node <sub>
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`,
			html.ElementNode, false},
		// missing sub-node <sub> and trees do not have the same structure
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`,
			html.TextNode, false},
		// additionnal sibling node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p><p class="ex2">outside</p>`,
			html.ElementNode,
			false},
		// additionnal sub-node <p class="ex2"
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub><p class="ex2">inside</p></p>`,
			html.ElementNode,
			false},
	}
	b := new(bytes.Buffer)
	fmt.Fprint(b, fragments[0].s)
	original, err := html.Parse(b)
	// Buffer is empty after parsing
	if err != nil {
		t.Errorf("%v while parsing %s", err, fragments[0].s)
	} else {
		// No Sibling in the document after parsing
		ExploreNode(original, "", original.Type)
		fmt.Println("---")
	}
	for i, f := range fragments {
		fmt.Fprint(b, f.s)
		n, err := html.Parse(b)
		if err != nil {
			t.Errorf("---FAIL %d: %v while parsing %s", i, err, f.s)
		}
		if r := IdenticalNodes(original, n, f.t); r != nil && f.equal {
			// if r := compareNodeRecursive(original, n); r != nil {
			t.Errorf("---FAIL(%d): %s differs from [%s]", i, f.s, PrintData(r))
			// exploreNode(r, "", html.ErrorNode)
		} else if r == nil && !f.equal {
			t.Errorf("---FAIL(%d): no difference found with %s", i, f.s)
		} else if r != nil && !f.equal {
			fmt.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
		}
	}
}
