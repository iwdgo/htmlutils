package parsing

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
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
	// Interverting attributes has not effect
	tag1.Attr = []html.Attribute{tag1.Attr[1], tag1.Attr[0]}
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	// Changing one value while keeping the number of attributes
	a := tag1.Attr[1]
	a.Val = "variable"
	tag1.Attr = []html.Attribute{tag1.Attr[1], a}
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found different\n", tag1, tag2)
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
	d, _ := ParseFile(f2)
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
			0, "caption", "",
			[]html.Attribute{}}, true},
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "p", "",
			[]html.Attribute{{"", "class", "not-found"}}}, false},
		{html.Node{nil, nil, nil, nil, nil, html.ElementNode,
			0, "p", "ns", // +0.9%
			[]html.Attribute{{"", "class", "not-found"}}}, false},
	}

	for _, m := range tagsToFind {
		want := PrintData(&m.n)
		p, _ := ParseFile(f)
		if o := FindNode(p, m.n); o == nil && m.f { // Should be found
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
		// Detected
		// identical to itself
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p>`,
			true},
		// missing letter in text
		{`<p class="ex1">HTML Fragment to compare against <em>other below</em> to test <sub>diffs</sub></p>`,
			false},
		// missing sub-node <sub>
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test diffs</p>`,
			false},
		// Not detected
		// additionnal sibling outside <p class="ex2">
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub></p><p class="ex2">outside</p>`,
			true},
		// additionnal sibling inside <p class="ex2">
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub><p class="ex2">inside</p></p>`,
			true},
		// TODO Missing closing tag
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub>`,
			true,
		},
		// Missing ending closing tags
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs`,
			true,
		},
	}
	b := new(bytes.Buffer)
	fmt.Fprint(b, fragments[0].s)
	original, err := html.Parse(b)
	// Buffer is empty after parsing
	if err != nil {
		t.Errorf("%v while parsing %s", err, fragments[0].s)
	} else {
		// No Sibling in the document after parsing
		// log.Printf("F %v\nL %v\nN %v\nP %v\n", original.FirstChild, original.LastChild, original.NextSibling, original.PrevSibling)
		ExploreNode(original, "", html.ErrorNode)
		log.Println("---")
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
			log.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
		} else {
			// Nothing is printed because the difference is not detected
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
		log.Println("---")
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
			log.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
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
		log.Println("---")
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
			log.Printf("---PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
		}
	}
}

func TestIsTextTag(t *testing.T) {
	titles := []struct {
		s string // want value
		b bool   // expected result
	}{
		{"Cayley table of Z<sub>3</sub> o Z<sub>3</sub> algebra of order 9", true},
		{"Wrong title on empty buffer", false},
	}
	tag := "caption"

	resp, err := http.Get("https://sitecloud-1266.appspot.com/displaytable?algebra=Z3%20o%20Z3")
	if err != nil {
		t.Error(err)
	}

	for _, ref := range titles {
		if err = IsTextTag(resp.Body, tag, ref.s); err != nil && ref.b && err.Error() != "findtag: tag not found" {
			t.Error(err)
		} else if err == nil && !ref.b {
			t.Error("titles are identical and should not")
		}
	}

	resp, err = http.Get("https://sitecloud-1266.appspot.com/displaytable?algebra=Z3%20o%20Z3")
	if err != nil {
		t.Error(err)
	}
	for _, ref := range titles[1:] {
		if err = IsTextTag(resp.Body, tag, ref.s); err == nil && !ref.b {
			t.Error("titles are identical and should not")
		}
	}
}

func TestIsTextNode(t *testing.T) {
	resp, err := http.Get("https://sitecloud-1266.appspot.com/displaytable?algebra=Z3%20o%20Z3")
	if err != nil {
		t.Error(err)
	}
	titles := []struct {
		s string // want value
		b bool   // expected result
	}{
		{"Cayley table of Z<sub>3</sub> o Z<sub>3</sub> algebra of order 9", true},
		{"Wrong title on empty buffer", false},
	}
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "caption"
	// TODO Add attributes
	// n.Attr = []html.Attribute{{"", "class", "fixed"}}
	for _, ref := range titles {
		if err = IsTextNode(resp.Body, &n, ref.s); err != nil && ref.b {
			t.Error(err)
		}
	}

}

func TestParseFileErrorFile(t *testing.T) {
	if _, err := ParseFile("doesnotexist"); !os.IsNotExist(err) {
		t.Error(err)
	}
}
