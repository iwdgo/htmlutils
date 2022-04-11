package parsing

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	f1 = "want_Test1.html"
	f2 = "want_Test2.html"
)

var parsingErr = []error{
	nil,
	errors.New("not found"),
	errors.New("not equal"),
	errors.New("attributes are differing"),
}

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
	tag2.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "fixed"}}
	if Equal(&tag1, &tag2) {
		t.Errorf("%v != %v are found equal\n", tag1, tag2)
	}
	tag1.Attr = tag2.Attr
	// .Attr[0] == .Attr[0]
	if !Equal(&tag1, &tag2) {
		t.Errorf("%v == %v are found different\n", tag1, tag2)
	}
	// .Attr[1] != .Attr[1]
	tag2.Attr = append(tag1.Attr, html.Attribute{Namespace: "", Key: "style", Val: "h2"})
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
	tagToFind.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex2"}}
	want := "p (Element) [{ class ex2}]"
	if s := PrintData(&tagToFind); s != want {
		t.Errorf("printData: got %s, want %s", s, want)
	}
}

func TestGetText(t *testing.T) {
	b := new(bytes.Buffer)
	_, _ = fmt.Fprint(b, HTMLf)
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
// account. Exact match between nodes is expected.
func TestFindNodes(t *testing.T) {
	f := f2
	tagsToFind := []struct {
		n html.Node
		f bool
	}{
		{html.Node{Type: html.ElementNode, Data: "table",
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "fixed"}}}, true},
		{html.Node{Type: html.ElementNode, Data: "caption"}, true},
		{html.Node{Type: html.ElementNode, Data: "p",
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "not-found"}}}, false},
		{html.Node{Type: html.ElementNode, Data: "p", Namespace: "ns", // +0.9%
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "ex1"}}}, false},
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
	n.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex2"}}
	if !Equal(IncludedNode(nil, &n), &n) {
		t.Errorf("nil does not include any node")
	}
}

// Using nil to include (+0.9%)
func TestIncludedNodeTyped(t *testing.T) {
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "p"
	n.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex2"}}
	if !Equal(IncludedNodeTyped(nil, &n, html.ErrorNode), &n) {
		t.Errorf("nil does not include any node")
	}
}

// Using nil to include (+0.8%)
func TestIdenticalNilNodes(t *testing.T) {
	var n html.Node
	n.Type = html.ElementNode
	n.Data = "p"
	n.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex2"}}
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
		// missing closing tag as accepted in HTML5
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs</sub>`,
			true,
		},
		// Missing ending closing tags
		{`<p class="ex1">HTML Fragment to compare against <em>others below</em> to test <sub>diffs`,
			true,
		},
	}
	b := new(bytes.Buffer)
	_, _ = fmt.Fprint(b, fragments[0].s)
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
		_, _ = fmt.Fprint(b, f.s)
		n, err := html.Parse(b)
		if err != nil {
			t.Errorf("---FAIL %d: %v while parsing %s", i, err, f.s)
		}
		if r := IncludedNode(original, n); r != nil && f.equal {
			// if r := compareNodeRecursive(original, n); r != nil {
			t.Errorf("--- FAIL(%d): %s differs from [%s]", i, f.s, PrintData(r))
			// exploreNode(r, "", html.ErrorNode)
		} else if r == nil && !f.equal {
			t.Errorf("--- FAIL(%d): no difference found with %s", i, f.s)
		} else if r != nil && !f.equal {
			log.Printf("=== PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
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
	_, _ = fmt.Fprint(b, fragments[0].s)
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
		_, _ = fmt.Fprint(b, f.s)
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
	_, _ = fmt.Fprint(b, fragments[0].s)
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
		_, _ = fmt.Fprint(b, f.s)
		n, err := html.Parse(b)
		if err != nil {
			t.Errorf("--- FAIL %d: %v while parsing %s", i, err, f.s)
		}
		if r := IdenticalNodes(original, n, f.t); r != nil && f.equal {
			// if r := compareNodeRecursive(original, n); r != nil {
			t.Errorf("--- FAIL(%d): %s differs from [%s]", i, f.s, PrintData(r))
			// exploreNode(r, "", html.ErrorNode)
		} else if r == nil && !f.equal {
			t.Errorf("--- FAIL(%d): no difference found with %s", i, f.s)
		} else if r != nil && !f.equal {
			log.Printf("=== PASS(%d): %s differs from: %s\n", i, f.s, PrintData(r))
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
	f, err := ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	b := new(bytes.Buffer)
	for _, ref := range titles {
		b.Reset()
		b.Write(f)
		err = IsTextTag(ioutil.NopCloser(b), tag, ref.s)
		if err == nil && ref.b {
			// success as expected
		} else if err != nil && !ref.b && strings.Contains(err.Error(), "not found") {
			t.Error(err)
		} else if err != nil && !ref.b && strings.Contains(err.Error(), "texts differ") {
			// failed as expected
		} else {
			t.Error("unknown error:", err)
		}
	}

	err = IsTextTag(ioutil.NopCloser(b), tag, titles[0].s)
	if err == nil {
		t.Error("titles are identical and should not")
	} else if strings.Contains(err.Error(), "not found") {
		// Not found as expected as b is empty
	} else {
		t.Errorf("got %v, want not found type error", err)
	}
}

func TestIsTextNode(t *testing.T) {
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
	// n.Attr = nil // empty array is invalid
	f, err := ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	b := new(bytes.Buffer)

	/* DEBUG
	dt, err := html.Parse(b)
	if err != nil {
		log.Println(err)
	}
	log.Println(FindTag(dt, tag, html.ElementNode))
	*/

	for _, ref := range titles {
		b.Reset()
		b.Write(f)
		err = IsTextNode(ioutil.NopCloser(b), &n, ref.s)
		if err == nil && ref.b {
			// success as expected
		} else if err != nil && !ref.b && strings.Contains(err.Error(), "not found") {
			t.Error(err)
		} else if err != nil && !ref.b && strings.Contains(err.Error(), "texts differ") {
			// failed as expected
		} else {
			t.Error("unknown error:", err)
		}
	}

	// On empty buffer
	err = IsTextNode(ioutil.NopCloser(b), &n, "does-not-matter")
	if err == nil {
		t.Fatal("node was unexpectedly found")
	}
	if want := "got findnode: node caption not found."; strings.Contains(err.Error(), want) {
		t.Errorf("error message: got %v, want %v", err, want)
	}

	// Node does not match on attributes
	b.Reset()
	b.Write(f)
	n.Attr = []html.Attribute{{Namespace: "", Key: "class", Val: "ex2"}}
	err = IsTextNode(ioutil.NopCloser(b), &n, "does-not-matter")
	if err == nil {
		t.Fatal("node was unexpectedly found")
	}
	if want := "got findnode: node caption not found."; strings.Contains(err.Error(), want) {
		t.Errorf("error message: got %v, want %v", err, want)
	}
}

func TestParseFileErrorFile(t *testing.T) {
	if _, err := ParseFile("doesnotexist"); !os.IsNotExist(err) {
		t.Error(err)
	}
}

// Testing the the search of a node based on values which are not pointers,
// i.e. the tree of the node is not taken into account.
func TestAttrIncluded(t *testing.T) {
	tagsToFind := []struct {
		n html.Node
		f error
	}{
		// On attribute missing
		{html.Node{Type: html.ElementNode, Data: "p",
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "ex1"}}}, nil},
		// No attribute
		{html.Node{Type: html.ElementNode, Data: "p"}, nil},
		// Wrong namespace
		{html.Node{Type: html.ElementNode, Data: "p", Namespace: "ns",
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "ex1"}}}, parsingErr[2]},
		// Wrong value of attribute
		{html.Node{Type: html.ElementNode, Data: "p",
			Attr: []html.Attribute{{Namespace: "", Key: "class", Val: "not-found"}}}, parsingErr[3]},
	}

	b := new(bytes.Buffer)
	for _, m := range tagsToFind {
		want := PrintData(&m.n)
		_, _ = fmt.Fprint(b, HTMLf2)
		p, err := html.Parse(b) // Any parsing error would occurred elsewhere
		if err != nil {
			t.Error(err)
		}
		o := FindTag(p, m.n.Data, html.ElementNode)
		switch m.f {
		case nil:
			if o == nil {
				t.Errorf("<%s> not found in %s\n", PrintData(&m.n), want)
			}
			continue
		case parsingErr[1]:
			if o == nil {
				// Not found as expected
			} else {
				t.Errorf("<%s> found as unexpected in %s\n", PrintData(&m.n), want)
			}
			continue
		case parsingErr[2]:
			if m.n.Namespace != o.Namespace {
				// Different as expected
			} else {
				t.Errorf("<%s> has invalid namespace in %s\n", PrintData(&m.n), want)
			}
			continue
		}

		inc := AttrIncluded(&m.n, o)
		if inc && m.f == parsingErr[3] {
			t.Errorf("<%s> attributes are included as unexpected in %s\n", PrintData(&m.n), want)
		}
		if !inc && m.f == nil {
			t.Errorf("<%s> attributes included as unexpected in %s\n", PrintData(&m.n), want)
		}
		b.Reset()
	}
}

func TestAttrIncludedEmpty(t *testing.T) {
	var m, n html.Node
	if !AttrIncluded(&m, &n) {
		t.Errorf("attrincluded: unexpected failure when no attributes are available")
	}
}

func TestFindTags(t *testing.T) {
	const tablesize = 81
	m, _ := ParseFile(f1)
	na := FindTags(m, "button", html.ElementNode)
	if len(na) != tablesize {
		t.Errorf("got %d, want %d", len(na), tablesize)
	}

	var p html.Node
	p.Type = html.ElementNode
	p.Data = "caption"
	p.Attr = []html.Attribute{{Namespace: "", Key: "title", Val: "Neutral"}}
	i := 0
	for _, n := range na {
		if AttrIncluded(n, &p) {
			i++
		}
	}
	const neutral = 9
	if i != neutral {
		t.Errorf("attributes: got %d, want %d", i, neutral)
	}
}

// html.Parse() does not return an error for its structure.
// To cover error processing, the Read method always returns an error.
type ErrReader struct{ Error error }

func (e *ErrReader) Read([]byte) (int, error) {
	return 0, e.Error
}

func TestIsTextNodeParseError(t *testing.T) {
	const parseError = "failing html.Parse()"
	err := IsTextNode(io.NopCloser(&ErrReader{errors.New(parseError)}), nil, "")
	if err == nil {
		t.Error("expected to fail")
	}
	want := fmt.Sprintf("parsing: %s", parseError)
	if err.Error() != want {
		t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
	}
}
