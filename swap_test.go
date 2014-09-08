package sgf

import (
	"fmt"
	"testing"
)

type SwapTester interface {
	Swap(a, b interface{})
	SwapAt(a, b int)
	Remove(a interface{})
	RemoveAt(a int)
	Len() int
	Object(i int) interface{}
	TestObjects() (a, b, c interface{})
	String(a interface{}) string
}

func TestNodePropertySwaps(t *testing.T) {
	testGenericSwapsAndRemoves(t, &NodeTester{})
}

func TestGameTreeNodeSwaps(t *testing.T) {
	testGenericSwapsAndRemoves(t, &GameTreeNodeTester{})
}

func TestGameTreeChildTreeSwaps(t *testing.T) {
	testGenericSwapsAndRemoves(t, &GameTreeChildTreeTester{})
}

func TestCollectionGameTreeSwaps(t *testing.T) {
	testGenericSwapsAndRemoves(t, &CollectionTester{})
}

func testGenericSwapsAndRemoves(t *testing.T, st SwapTester) {
	validateOrder := func(desc string, st SwapTester, objects ...interface{}) {
		if st.Len() != len(objects) {
			t.Fatalf("%s. Object len mismatch. Was %d, expected %d.", desc, st.Len(), len(objects))
		}

		for i, o := range objects {
			if st.Object(i) != o {
				t.Fatalf("Object[%d] mismatch. Was %s, expected: %s", i, st.String(st.Object(i)), st.String(o))
			}
		}
	}

	p1, p2, p3 := st.TestObjects()

	// Check initial order
	validateOrder("check initial order", st, p1, p2, p3)

	// Swap 1. and 3. element
	st.Swap(p1, p3)
	validateOrder("Swap 1. and 3. element", st, p3, p2, p1)

	// Remove 2. element
	st.Remove(p2)
	validateOrder("Remove 2. element", st, p3, p1)

	// Same thing with indices
	st = &NodeTester{nil}
	p1, p2, p3 = st.TestObjects()

	// Check initial order
	validateOrder("check initial order", st, p1, p2, p3)

	// Swap 1. and 3. element
	st.SwapAt(0, 2)
	validateOrder("Swap 1. and 3. element", st, p3, p2, p1)

	// Remove 2. element
	st.RemoveAt(1)
	validateOrder("Remove 2. element", st, p3, p1)
}

func TestSwapRemoveInvalidIndex(t *testing.T) {
	c, gt, n := NewCollection()
	p := n.NewProperty("FF", "1")

	nonexistentGameTree := &GameTree{}
	nonexistentNode := &Node{}
	nonexistentProperty := &Property{}

	var tests = []struct {
		desc string
		fn   func()
	}{
		{"c.SwapGameTreesAt(0, 1)", func() { c.SwapGameTreesAt(0, 1) }},
		{"c.SwapGameTrees(gt, nonexistentGameTree)", func() { c.SwapGameTrees(gt, nonexistentGameTree) }},
		{"gt.SwapGameTreesAt(0, 1)", func() { gt.SwapGameTreesAt(0, 1) }},
		{"gt.SwapGameTrees(gt, nonexistentGameTree)", func() { gt.SwapGameTrees(gt, nonexistentGameTree) }},
		{"gt.SwapNodesAt(0, 1)", func() { gt.SwapNodesAt(0, 1) }},
		{"gt.SwapNodes(n, 1)", func() { gt.SwapNodes(n, nonexistentNode) }},
		{"n.SwapPropertiesAt(0, 1)", func() { n.SwapPropertiesAt(0, 1) }},
		{"n.SwapProperties(p, nonexistentProperty)", func() { n.SwapProperties(p, nonexistentProperty) }},
		{"c.RemoveGameTreeAt(1)", func() { c.RemoveGameTreeAt(1) }},
		{"c.RemoveGameTree(nonexistentGameTree)", func() { c.RemoveGameTree(nonexistentGameTree) }},
		{"gt.RemoveGameTreeAt(1)", func() { gt.RemoveGameTreeAt(1) }},
		{"gt.RemoveGameTree(nonexistentGameTree)", func() { gt.RemoveGameTree(nonexistentGameTree) }},
		{"gt.RemoveNodeAt(1)", func() { gt.RemoveNodeAt(1) }},
		{"gt.RemoveNode(nonexistentNode)", func() { gt.RemoveNode(nonexistentNode) }},
	}

	for _, test := range tests {

		panicked := fnPanics(test.fn)

		if !panicked {
			t.Errorf("%s should have panicked.", test.desc)
		}
	}
}

//
// Tester for Node's Properties
//

type NodeTester struct {
	node *Node
}

func (t NodeTester) Swap(a, b interface{}) {
	t.node.SwapProperties(a.(*Property), b.(*Property))
}

func (t NodeTester) SwapAt(a, b int) {
	t.node.SwapPropertiesAt(a, b)
}

func (t NodeTester) Remove(a interface{}) {
	t.node.RemoveProperty(a.(*Property))
}

func (t NodeTester) RemoveAt(a int) {
	t.node.RemovePropertyAt(a)
}

func (t NodeTester) Len() int {
	return len(t.node.Properties)
}

func (t NodeTester) Object(i int) interface{} {
	return t.node.Properties[i]
}

func (t *NodeTester) TestObjects() (a, b, c interface{}) {
	_, _, t.node = NewCollection()
	a = t.node.NewProperty("FF", "1")
	b = t.node.NewProperty("FF", "3")
	c = t.node.NewProperty("FF", "5")
	return
}

func (t NodeTester) String(a interface{}) string {
	return fmt.Sprintf("%s[%s]", a.(*Property).Ident, a.(*Property).Values[0])
}

//
// Tester for GameTree's Nodes
//

type GameTreeNodeTester struct {
	gt *GameTree
}

func (t GameTreeNodeTester) Swap(a, b interface{}) {
	t.gt.SwapNodes(a.(*Node), b.(*Node))
}

func (t GameTreeNodeTester) SwapAt(a, b int) {
	t.gt.SwapNodesAt(a, b)
}

func (t GameTreeNodeTester) Remove(a interface{}) {
	t.gt.RemoveNode(a.(*Node))
}

func (t GameTreeNodeTester) RemoveAt(a int) {
	t.gt.RemoveNodeAt(a)
}

func (t GameTreeNodeTester) Len() int {
	return len(t.gt.Nodes)
}

func (t GameTreeNodeTester) Object(i int) interface{} {
	return t.gt.Nodes[i]
}

func (t *GameTreeNodeTester) TestObjects() (a, b, c interface{}) {
	_, gt, n1 := NewCollection()
	n2 := gt.NewNode()
	n3 := gt.NewNode()

	n1.NewProperty("FF", "1")
	n2.NewProperty("FF", "3")
	n3.NewProperty("FF", "5")

	t.gt = gt
	return n1, n2, n3
}

func (t GameTreeNodeTester) String(a interface{}) string {
	return fmt.Sprintf(";%s[%s]", a.(*Node).Properties[0].Ident, a.(*Node).Properties[0].Values[0])
}

//
// Tester for GameTree's child GameTrees
//

type GameTreeChildTreeTester struct {
	gt *GameTree
}

func (t GameTreeChildTreeTester) Swap(a, b interface{}) {
	t.gt.SwapGameTrees(a.(*GameTree), b.(*GameTree))
}

func (t GameTreeChildTreeTester) SwapAt(a, b int) {
	t.gt.SwapGameTreesAt(a, b)
}

func (t GameTreeChildTreeTester) Remove(a interface{}) {
	t.gt.RemoveGameTree(a.(*GameTree))
}

func (t GameTreeChildTreeTester) RemoveAt(a int) {
	t.gt.RemoveGameTreeAt(a)
}

func (t GameTreeChildTreeTester) Len() int {
	return len(t.gt.GameTrees)
}

func (t GameTreeChildTreeTester) Object(i int) interface{} {
	return t.gt.GameTrees[i]
}

func (t *GameTreeChildTreeTester) TestObjects() (a, b, c interface{}) {
	_, gt, _ := NewCollection()

	gt1, n1 := gt.NewGameTree()
	gt2, n2 := gt.NewGameTree()
	gt3, n3 := gt.NewGameTree()

	n1.NewProperty("FF", "1")
	n2.NewProperty("FF", "3")
	n3.NewProperty("FF", "5")

	t.gt = gt
	return gt1, gt2, gt3
}

func (t GameTreeChildTreeTester) String(a interface{}) string {
	return fmt.Sprintf("(;%s[%s])", a.(*GameTree).Nodes[0].Properties[0].Ident, a.(*GameTree).Nodes[0].Properties[0].Values[0])
}

//
// Tester for Collection's GameTrees
//

type CollectionTester struct {
	c *Collection
}

func (t CollectionTester) Swap(a, b interface{}) {
	t.c.SwapGameTrees(a.(*GameTree), b.(*GameTree))
}

func (t CollectionTester) SwapAt(a, b int) {
	t.c.SwapGameTreesAt(a, b)
}

func (t CollectionTester) Remove(a interface{}) {
	t.c.RemoveGameTree(a.(*GameTree))
}

func (t CollectionTester) RemoveAt(a int) {
	t.c.RemoveGameTreeAt(a)
}

func (t CollectionTester) Len() int {
	return len(t.c.GameTrees)
}

func (t CollectionTester) Object(i int) interface{} {
	return t.c.GameTrees[i]
}

func (t *CollectionTester) TestObjects() (a, b, c interface{}) {
	coll, gt1, n1 := NewCollection()
	gt2, n2 := coll.NewGameTree()
	gt3, n3 := coll.NewGameTree()

	n1.NewProperty("FF", "1")
	n2.NewProperty("FF", "3")
	n3.NewProperty("FF", "5")

	t.c = coll
	return gt1, gt2, gt3
}

func (t CollectionTester) String(a interface{}) string {
	return fmt.Sprintf("(;%s[%s])", a.(*GameTree).Nodes[0].Properties[0].Ident, a.(*GameTree).Nodes[0].Properties[0].Values[0])
}
