package sgf

import (
	"bytes"
	"io/ioutil"
)

//
// Collection
//

// Create a new collection. Collection always contains at least one GameTree and one Node and those are also
// automatically created.
func NewCollection() (*Collection, *GameTree, *Node) {
	collection := &Collection{}
	gameTree, node := collection.NewGameTree()

	return collection, gameTree, node
}

// Creates a new GameTree (and Node) and add those to the collection.
func (collection *Collection) NewGameTree() (*GameTree, *Node) {
	gameTree := &GameTree{}
	collection.AddGameTree(gameTree)
	node := gameTree.NewNode()
	return gameTree, node
}

// Adds a new GameTree to this Collection.
func (collection *Collection) AddGameTree(gameTree *GameTree) {
	collection.GameTrees = append(collection.GameTrees, gameTree)
}

// Swaps the order of the given GameTrees in this Collection.
func (collection *Collection) SwapGameTrees(gt1, gt2 *GameTree) {
	collection.SwapGameTreesAt(collection.gameTreeIndex(gt1), collection.gameTreeIndex(gt2))
}

// Swaps the order of GameTrees at the given indices in this Collection.
func (collection *Collection) SwapGameTreesAt(i, j int) {
	collection.GameTrees[i], collection.GameTrees[j] = collection.GameTrees[j], collection.GameTrees[i]
}

// Removes the given GameTree from this Collection.
func (collection *Collection) RemoveGameTree(gameTree *GameTree) {
	collection.RemoveGameTreeAt(collection.gameTreeIndex(gameTree))
}

// Removes GameTree from this Collection at the given index.
func (collection *Collection) RemoveGameTreeAt(i int) {
	collection.GameTrees = append(collection.GameTrees[:i], collection.GameTrees[i+1])
}

// returns the index of the given GameTree in this Collection or -1 if GameTree is not present.
func (collection *Collection) gameTreeIndex(gameTree *GameTree) int {
	for i, gt := range collection.GameTrees {
		if gt == gameTree {
			return i
		}
	}

	return -1
}

//
// GameTree
//

// Creates a new Node and add it to the GameTree
func (gameTree *GameTree) NewNode() *Node {
	node := &Node{}
	gameTree.AddNode(node)
	return node
}

// Adds a new Node to this GameTree.
func (gameTree *GameTree) AddNode(node *Node) {
	gameTree.Nodes = append(gameTree.Nodes, node)
}

// Creates a new GameTree (and Node) and add it as a child.
func (gameTree *GameTree) NewGameTree() (*GameTree, *Node) {
	newGameTree := &GameTree{}
	gameTree.AddGameTree(newGameTree)
	node := newGameTree.NewNode()
	return newGameTree, node
}

// Adds a new child GameTree to this GameTree.
func (gameTree *GameTree) AddGameTree(newGameTree *GameTree) {
	gameTree.GameTrees = append(gameTree.GameTrees, newGameTree)
}

// Swaps the order of the given child GameTrees in this GameTree.
func (gameTree *GameTree) SwapGameTrees(gt1, gt2 *GameTree) {
	gameTree.SwapGameTreesAt(gameTree.gameTreeIndex(gt1), gameTree.gameTreeIndex(gt2))
}

// Swaps the oder of child GameTrees at the given indices in this GameTree.
func (gameTree *GameTree) SwapGameTreesAt(i, j int) {
	gameTree.GameTrees[i], gameTree.GameTrees[j] = gameTree.GameTrees[j], gameTree.GameTrees[i]
}

// Swaps the order of the given Nodes in this GameTree.
func (gameTree *GameTree) SwapNodes(n1, n2 *Node) {
	gameTree.SwapNodesAt(gameTree.nodeIndex(n1), gameTree.nodeIndex(n2))
}

// Swaps the order of Nodes at the given indices in this GameTree.
func (gameTree *GameTree) SwapNodesAt(i, j int) {
	gameTree.Nodes[i], gameTree.Nodes[j] = gameTree.Nodes[j], gameTree.Nodes[i]
}

// Removes the given child GameTree from this GameTree.
func (gameTree *GameTree) RemoveGameTree(gt *GameTree) {
	gameTree.RemoveGameTreeAt(gameTree.gameTreeIndex(gt))
}

// Removes child GameTree from this Collection at the given index.
func (gameTree *GameTree) RemoveGameTreeAt(i int) {
	gameTree.GameTrees = append(gameTree.GameTrees[:i], gameTree.GameTrees[i+1])
}

// Removes the given Node from this GameTree.
func (gameTree *GameTree) RemoveNode(n *Node) {
	gameTree.RemoveNodeAt(gameTree.nodeIndex(n))
}

// Removes Node from this Collection at the given index.
func (gameTree *GameTree) RemoveNodeAt(i int) {
	gameTree.Nodes = append(gameTree.Nodes[:i], gameTree.Nodes[i+1])
}

// returns the index of the given child GameTree in this GameTree or -1 if child GameTree is not present.
func (gameTree *GameTree) gameTreeIndex(childGameTree *GameTree) int {
	for i, gt := range gameTree.GameTrees {
		if gt == childGameTree {
			return i
		}
	}

	return -1
}

// returns the index of the given Node in this GameTree or -1 if Node is not present.
func (gameTree *GameTree) nodeIndex(node *Node) int {
	for i, n := range gameTree.Nodes {
		if n == node {
			return i
		}
	}

	return -1
}

//
// Node
//

// Creates a new Property and add it to the Node.
func (node *Node) NewProperty(ident string, values ...string) *Property {
	property := &Property{ident, values}
	node.AddProperty(property)
	return property
}

// Removes the given Property from this Node.
func (node *Node) RemoveProperty(property *Property) {
	for i, p := range node.Properties {
		if p == property {
			node.RemovePropertyAt(i)
			return
		}
	}
}

// Removes Property from this Node at the given index.
func (node *Node) RemovePropertyAt(i int) {
	node.Properties = append(node.Properties[:i], node.Properties[i+1:]...)
}

// Adds a new property to the Node.
func (node *Node) AddProperty(property *Property) {
	node.Properties = append(node.Properties, property)
}

// Swaps the order of the given Properties in this Node.
func (node *Node) SwapProperties(p1, p2 *Property) {
	node.SwapPropertiesAt(node.propertyIndex(p1), node.propertyIndex(p2))
}

// Swaps the order of Properties at the given indices in this Node.
func (node *Node) SwapPropertiesAt(i, j int) {
	node.Properties[i], node.Properties[j] = node.Properties[j], node.Properties[i]
}

// returns the index of the given Property in this Node or -1 if Property is not present.
func (node *Node) propertyIndex(property *Property) int {
	for i, p := range node.Properties {
		if p == property {
			return i
		}
	}

	return -1
}

//
// Other
//

// Check if the collection is valid.
// Collection is not valid if:
//   - Collection does not have any GameTrees.
//   - Any GameTree inside the collection does not have any Nodes.
//   - Any Property inside the collection does not have any values.
func (collection *Collection) Valid() bool {
	// Collection must contain at least one GameTree
	if len(collection.GameTrees) == 0 {
		return false
	}

	// Check GameTrees
	for _, gameTree := range collection.GameTrees {
		if !validGameTree(gameTree) {
			return false
		}
	}

	return true
}

func validGameTree(gameTree *GameTree) bool {
	// GameTree must contain at least one Node
	if len(gameTree.Nodes) == 0 {
		return false
	}

	// Check Properties
	for _, node := range gameTree.Nodes {
		for _, property := range node.Properties {
			// Check Ident
			if property.Ident == "" {
				return false
			}

			// Check Values
			if len(property.Values) == 0 {
				return false
			}
		}
	}

	// Check child GameTrees recursively
	for _, childGameTree := range gameTree.GameTrees {
		if !validGameTree(childGameTree) {
			return false
		}
	}

	return true
}

// Used to format SGF format.
type SgfFormat struct {
	NewLineBetweenGameTrees bool // put each gameTree to its own line
	NewLineAlsoBetweenNodes bool // also put each node to its own line, only works if NewLineBetweenGameTrees = true
	IndentationLevel        int  // how many whitespaces are used when indenting
}

var (
	DefaultSgfFormat    = SgfFormat{true, true, 4}
	NoNewLinesSgfFormat = SgfFormat{false, false, 0}
)

// Converts the collection to SGF format.
func (collection *Collection) Sgf(format SgfFormat) string {
	var buffer bytes.Buffer

	if !collection.Valid() {
		panic("collection is not valid.")
	}

	for _, gameTree := range collection.GameTrees {
		appendGameTree(&buffer, gameTree, format, 0)
	}

	return buffer.String()
}

func appendGameTree(buffer *bytes.Buffer, gameTree *GameTree, format SgfFormat, level int) {
	// Add newlines between GameTrees if required
	if format.NewLineBetweenGameTrees && level > 0 {
		buffer.WriteRune('\n')

		for i := 0; i < level; i++ {
			for j := 0; j < format.IndentationLevel; j++ {
				buffer.WriteRune(' ')
			}
		}
	}

	// Start of GameTree
	buffer.WriteRune('(')

	for i, node := range gameTree.Nodes {
		// Add new lines between nodes if required
		if format.NewLineBetweenGameTrees && format.NewLineAlsoBetweenNodes && i > 0 {
			buffer.WriteRune('\n')

			for i := 0; i < level; i++ {
				for j := 0; j < format.IndentationLevel; j++ {
					buffer.WriteRune(' ')
				}
			}

			// Add one more whitespace to line up nodes. Otherwise we get lines like:
			// (;FF[4]
			// ;SZ[19]
			// Instead of:
			// (;FF[4]
			//  ;SZ[19]
			buffer.WriteRune(' ')
		}

		// Start of Node
		buffer.WriteRune(';')

		for _, property := range node.Properties {
			// Property ident
			buffer.WriteString(property.Ident)

			// Property values
			for _, value := range property.Values {
				buffer.WriteRune('[')
				buffer.WriteString(value)
				buffer.WriteRune(']')
			}
		}
	}

	// Child GameTrees
	for _, childGameTree := range gameTree.GameTrees {
		appendGameTree(buffer, childGameTree, format, level+1)
	}

	// End of GameTree
	buffer.WriteRune(')')
}

var readFileFunc func(string) ([]byte, error) = ioutil.ReadFile

// Parse given filename as a SGF file.
func ParseSgfFile(filename string) (*Collection, error) {
	bytes, err := readFileFunc(filename)
	if err != nil {
		return nil, err
	}

	return ParseSgf(string(bytes))
}

// Parse given data as a SGF file.
func ParseSgf(data string) (*Collection, error) {
	lexemes, err := lexicalAnalysis(data)
	if err != nil {
		return nil, err
	}

	collection, err := parse(lexemes)

	return collection, err
}
