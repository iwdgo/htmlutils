// Package parsing provides basic search and comparison of HTML documents.
// To limit storage of references, it uses the net/html package and its Node type to structure HTML.
// Search a tag in a Node with options
// - searching a tag based on its name whatever attributes where its type is optional
// - searching a tag based on its non-pointer values: type, name, attribute and namespace
// - comparing tags including list of attributes where order is irrelevant
// - comparing Node structures with an optional type
// Three ways to print a node tree.
// - Select one type and stop on one value
// - ElementNode or all
// - Complete with indentation
// Good to know:
// - a non-matching closed tag is one element.
// - a non-closed tag is closed by the following opening tag.
//   The elements that follow are discarded as the tag is closed by the parser.
package parsing

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

var nodeTypeNames = []string{"Error", "Text", "Document", "Element", "Comment", "DocType"}

/*
HTML structure requires to locate the tag from which comparison starts and compare all tags until the comparison ends
All Print and Search funcs are recursive. Every Sibling of the FirstChild is searched for the same value.
Tests demonstrate usage of the library

TODO Siblings might not have the same order
Problem is the order is sometimes relevant. So relaxing completely the order of the siblings might not be useful.
TODO iota (html.ErrorNode) seems difficult to produce and looks like the right default
You can first check that no ErrorNode was created before comparing for instance
*/
// ParseFile returns a *Node containing the parsed file
// Func cannot be tested and maximum test coverage is 98.3%
func ParseFile(f string) *html.Node {
	file, err := os.Open(f)
	if err != nil {
		log.Fatalln("opening file failed", err)
	}
	defer file.Close()
	n, err := html.Parse(file)
	if err != nil {
		log.Fatalln("parsing failed:", err)
	}
	return n
}

// String-based search

// ExploreNode prints node tags with name s and type t
// Without name, all tags are printed
// When type ErrorNode (iota == 0) prints tags of all types
func ExploreNode(n *html.Node, s string, t html.NodeType) {
	if n.Type == t || t == html.ErrorNode {
		if n.Data == s || s == "" {
			fmt.Printf(" %s", printData(n))
		}
	}
	// Something will print
	if n.FirstChild != nil {
		fmt.Print("\n") // Siblings on one line
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ExploreNode(c, s, t)
	}
}

// PrintTags prints node structure until a tag name is found (whatever attributes)
// Without name, all tags are printed
// tagOnly selects ElementNode, otherwise tags are printed whatever type.
// If node tree has no Errornode, there is no difference with previous
// i.e. exploreNode(n, "", html.ErrorNode) prints nothing then both are equivalent.
func PrintTags(n *html.Node, s string, tagOnly bool) {
	if tagOnly && n.Type == html.ElementNode { // tag is true and only tags are dumped
		fmt.Println(printData(n))
	} else if !tagOnly { // Otherwise, all nodes
		fmt.Println(printData(n))
	}
	if s != "" && n.Type == html.ElementNode && n.Data == s {
		fmt.Printf("[%s] found. Stopping exploration\n", n.Data)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		PrintTags(c, s, tagOnly)
	}
}

// FindTag finds the first occurrence of a tag name (i.e. whatever its attributes).
// If ErrorNode is passed, any tag type will be searched
func FindTag(n *html.Node, s string, t html.NodeType) *html.Node {
	if n.Data == s && (n.Type == t || t == html.ErrorNode) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if f := FindTag(c, s, t); f != nil {
			return f
		}
	}
	return nil
}

// Search using Node
//
func indent(i int) (s string) {
	for j := 0; j < i; j++ {
		s += "."
	}
	return
}

// PrintNodes prints the tree structure of node m until n node is equal.
// If nil is passed, the complete node is printed.
// Values are indented based on the recursion depth d which is usually 0 when called
// html.ErrorNode (iota) displays every tag except the error node.
func PrintNodes(m, n *html.Node, t html.NodeType, d int) {
	if Equal(m, n) {
		fmt.Printf("\ntag found: %s", printData(m))
	}
	if m.FirstChild != nil {
		fmt.Printf("\n%s", indent(d)) // Siblings on one line
	}
	d++
	for o := m.FirstChild; o != nil; o = o.NextSibling {
		if o.Type == t || t == html.ErrorNode {
			fmt.Printf(" %s", printData(o))
		}
		PrintNodes(o, n, t, d)
	}
}

// GetText prints the text content of a tree structure like PrintNodes w/o any formatting
func GetText(m *html.Node, b *bytes.Buffer) {
	for o := m.FirstChild; o != nil; o = o.NextSibling {
		if o.Type == html.TextNode {
			_, _ = fmt.Fprint(b, o.Data)
		}
		GetText(o, b)
	}
}

// findAttr locates an attribute in a list of attributes
func findAttr(a html.Attribute, l []html.Attribute) bool {
	for _, d := range l {
		if a == d {
			return true
		}
	}
	return false
}

// attrEqual return true if list of attributes are equal whatever their order
func attrEqual(m, n *html.Node) bool {
	if len(m.Attr) == 0 && len(n.Attr) == 0 {
		return true
	}
	identicalAttr := true
	i := 0
	for identicalAttr && i < len(m.Attr) && i < len(n.Attr) {
		identicalAttr = identicalAttr && findAttr(m.Attr[i], n.Attr) // m.Attr[i] == n.Attr[i]
		i++
	}
	// i was incremented when one attribute was found and must have the length of each array
	if identicalAttr && i == len(m.Attr) && i == len(n.Attr) {
		return true
	}
	return false
}

// Equal returns true if all fields of nodes m and n are equal except pointers
// reflect.DeepEqual(tag1, tag2) is unusable as pointers are checked too
func Equal(m, n *html.Node) bool {
	// This test is something like reflect.TypeOf(m) == reflect.TypeOf(n)
	if m == nil && n == nil { // Passing untyped value panics otherwise
		return true
	} else if m == nil || n == nil {
		return false
	}
	return m.Type == n.Type && m.Data == n.Data && attrEqual(m, n) && m.Namespace == n.Namespace
}

// printData returns a string with Node information (not its relationships)
// nil will panic
func printData(n *html.Node) string {
	ns := ""
	if n.Namespace != "" {
		ns = " ns:[" + n.Namespace + "]"
	}
	nattr := ""
	if len(n.Attr) != 0 {
		nattr = fmt.Sprintf("%v", n.Attr)
	}
	return strings.TrimSpace(n.Data + " (" + nodeTypeNames[n.Type] + ") " + nattr + ns)
}

// FindNode find the first occurrence of a node
func FindNode(m *html.Node, n html.Node) *html.Node {
	if Equal(m, &n) {
		//
		return m
	}
	for c := m.FirstChild; c != nil; c = c.NextSibling {
		if f := FindNode(c, n); f != nil {
			return f
		}
	}
	// else keep searching by returning nil
	return nil
}

// Tree handling

// IncludeNode checks if n is included in m.
// Included means that the subtree is identical to m including order of siblings.
// If it is, nil is returned. Otherwise, the tag from which trees diverge is returned.
// If m has more tags than n, nil is returned as the search stops when one subtree exploration is exhausted.
func IncludedNode(m, n *html.Node) *html.Node {
	if !Equal(m, n) {
		// Return the non-nil value
		if m == nil {
			return n
		}
		return m // returning the tree that includes
	}
	// Looping over siblings of FirstChild
	nf := n.FirstChild
	for c := m.FirstChild; c != nil; c = c.NextSibling {
		// and comparing to the other tree in the same order
		if cn := IncludedNode(c, nf); cn != nil { // Some diff found - printing non-nil
			//fmt.Printf("cn (where different):%s\t", printData(cn))
			/* TODO Test is useless
			if c != nil {
				fmt.Printf("m child:%s\t", printData(c))
			}
			if nf != nil {
				fmt.Printf("n child:%s\n", printData(nf))
			} else {
				fmt.Println("n child (nil):", nf)
			}
			*/
			return cn
		}
		nf = nf.NextSibling
	}
	return nil
}

// IncludedNodeTyped is like IncludeNode where only tags of type t are compared
func IncludedNodeTyped(m, n *html.Node, t html.NodeType) *html.Node {
	if !Equal(m, n) {
		// Difference matters only if type is as requested
		// Returning the eventual non-nil value
		if m == nil && n != nil {
			return n
		} else if n == nil && m != nil {
			return m
		}
		if m.Type == t && n.Type == t {
			return m // returning the including node
		}
	}
	// Looping over siblings of FirstChild
	nf := n.FirstChild
	for c := m.FirstChild; c != nil; c = c.NextSibling {
		// and comparing to the other tree in the same order
		if cn := IncludedNodeTyped(c, nf, t); cn != nil { // Some diff found - printing non-nil
			if cn.Type == t {
				return cn
			}
		}
		nf = nf.NextSibling
	}
	return nil
}

// IdenticalNodes fails if trees have different size
func IdenticalNodes(m, n *html.Node, t html.NodeType) *html.Node {
	if !Equal(m, n) {
		// Difference matters only if type is as requested
		// Returning the eventual non-nil value
		if m == nil && n != nil {
			return n
		} else if n == nil && m != nil {
			return m
		}
		if m.Type == t && n.Type == t {
			return m // returning the including node
		}
	}
	// Looping over siblings of FirstChild
	nf := n.FirstChild
	for c := m.FirstChild; c != nil; c = c.NextSibling {
		// and comparing to the other tree in the same order
		if cn := IdenticalNodes(c, nf, t); cn != nil { // Some diff found - printing non-nil
			if cn.Type == t {
				return cn
			}
		}
		nf = nf.NextSibling
	}
	if nf != nil {
		return nf
	}
	return nil
}
