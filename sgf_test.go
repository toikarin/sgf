package sgf

import (
	"errors"
	"testing"
)

func TestReadFileError(t *testing.T) {
	oldReadFileFunc := readFileFunc
	defer func() {
		readFileFunc = oldReadFileFunc
	}()

	readFileFunc = func(filename string) ([]byte, error) {
		return nil, errors.New("")
	}

	_, err := ParseSgfFile("foo")

	if err == nil {
		t.Errorf("ParseSgfFile did not return error.")
	}
}

func TestReadFileOk(t *testing.T) {
	oldReadFileFunc := readFileFunc
	defer func() {
		readFileFunc = oldReadFileFunc
	}()

	readFileFunc = func(filename string) ([]byte, error) {
		return []byte("(;)"), nil
	}

	_, err := ParseSgfFile("foo")

	if err != nil {
		t.Errorf("ParseSgfFile returned error.")
	}
}

func TestSgf(t *testing.T) {
	var okTests = []struct {
		data   string
		wanted string
		format SgfFormat
	}{
		// empty
		{"(;)", "", NoNewLinesSgfFormat},
		{"(;FF[4]GM[1])", "", NoNewLinesSgfFormat},
		{"(;FF[4];GM[1](;PB[Black]))", "", NoNewLinesSgfFormat},
		{"(;FF[4]\n ;GM[1]\n    (;PB[Black]\n     ;PW[White]))", "", DefaultSgfFormat},
		{"(;FF[4]GM[1]C[username [rank\\]: \\\\o])", "", DefaultSgfFormat},
	}

	for _, test := range okTests {
		collection, err := ParseSgf(test.data)
		if err != nil {
			t.Errorf("ParseSgf(%s) returned error.", test.data)
			continue
		}

		sgf := collection.Sgf(test.format)
		wanted := test.wanted

		if test.wanted == "" {
			wanted = test.data
		}

		if wanted != sgf {
			t.Errorf("collection.Sgf(%s) mismatch. Got: %s.", test.data, sgf)
		}
	}
}

func TestSgfPanics(t *testing.T) {
	c := Collection{}

	panicked := fnPanics(func() { c.Sgf(DefaultSgfFormat) })

	if !panicked {
		t.Errorf("Sgf() should have panicked")
	}
}

func TestCollectionValid(t *testing.T) {
	var tests = []struct {
		collection Collection
		desc       string
	}{
		{Collection{}, "empty collection (missing GameTrees)"},
		{Collection{[]*GameTree{&GameTree{}}}, "GameTree missing nodes"},
		{Collection{[]*GameTree{&GameTree{[]*Node{&Node{[]*Property{&Property{"FF", []string{}}}}}, []*GameTree{}}}}, "Property missing value"},
		{Collection{[]*GameTree{&GameTree{[]*Node{&Node{[]*Property{&Property{"", []string{"4"}}}}}, []*GameTree{}}}}, "Property missing ident"},
		{Collection{[]*GameTree{&GameTree{[]*Node{&Node{}}, []*GameTree{&GameTree{}}}}}, "Child GameTree missing nodes"},
	}

	for _, test := range tests {
		if test.collection.Valid() {
			t.Errorf("Collection.Valid() failed on test: %s", test.desc)
		}
	}
}

func TestCollectionCreation(t *testing.T) {
	c, gt, n := NewCollection()

	// Check collection
	if len(c.GameTrees) != 1 {
		t.Fatalf("NewCollection(): New collection does not contain 1 gameTree. instead it has: %d", len(c.GameTrees))
	}

	if c.GameTrees[0] != gt {
		t.Fatalf("NewCollection(): returned GameTree does not much one in the collection.")
	}

	// Check GameTree
	if len(gt.Nodes) != 1 {
		t.Fatalf("NewCollection(): New GameTree does not contain 1 Node. instead it has: %d", len(gt.Nodes))
	}

	if len(gt.GameTrees) != 0 {
		t.Fatalf("NewCollection(): New GameTree does not contain 0 GameTree. instead it has: %d", len(gt.GameTrees))
	}

	if gt.Nodes[0] != n {
		t.Fatalf("NewCollection(): returned Node does not much one in the collection.")
	}

	// Check node
	if len(n.Properties) != 0 {
		t.Fatalf("NewCollection(): New Node does not contain 0 Properties. instead it has: %d", len(n.Properties))
	}

	// Add new property
	p := n.NewProperty("FF", "1")

	// Check node
	if len(n.Properties) != 1 {
		t.Fatalf("NewProperty(): Node does not contain 1 Properties. instead it has: %d", len(n.Properties))
	}

	if n.Properties[0] != p {
		t.Fatalf("NewProperty(): returned Property does not much one in the collection.")
	}

	// Check property
	if p.Ident != "FF" {
		t.Fatalf("NewProperty(): Ident mismatch. was: %d", n.Properties[0].Ident)
	}

	if len(p.Values) != 1 {
		t.Fatalf("NewProperty(): Property does not contain 1 value. instead it has: %d", len(p.Values))
	}

	// Add new GameTree
	gt2, n2 := gt.NewGameTree()

	// Check GameTree
	if len(gt.Nodes) != 1 {
		t.Fatalf("NewGameTree(): New GameTree does not contain 1 Node. instead it has: %d", len(gt.Nodes))
	}

	if len(gt.GameTrees) != 1 {
		t.Fatalf("NewGameTree(): New GameTree does not contain 1 GameTree. instead it has: %d", len(gt.GameTrees))
	}

	if gt.GameTrees[0] != gt2 {
		t.Fatalf("NewGameTree(): returned GameTree does not much one in the collection.")
	}

	// Check Child GameTree
	if len(gt2.Nodes) != 1 {
		t.Fatalf("NewGameTree(): New GameTree does not contain 1 Node. instead it has: %d", len(gt2.Nodes))
	}

	if len(gt2.GameTrees) != 0 {
		t.Fatalf("NewGameTree(): New GameTree does not contain 0 GameTree. instead it has: %d", len(gt2.GameTrees))
	}

	if gt2.Nodes[0] != n2 {
		t.Fatalf("NewGameTree(): returned Node does not much one in the collection.")
	}
}
