package sgf_test

import (
	"fmt"
	"github.com/blakki/sgf"
)

func ExampleParseSgfFile() {
	collection, err := sgf.ParseSgf("(;FF[4]GM[1])")
	if err != nil {
		panic(err)
	}

	for _, gameTree := range collection.GameTrees {
		for _, node := range gameTree.Nodes {
			for _, property := range node.Properties {
				fmt.Print(property.Ident)

				for _, value := range property.Values {
					fmt.Print(" ", value)
				}

				fmt.Println()
			}
		}
	}
	// Output: FF 4
	// GM 1
}

func ExampleCollection_Sgf() {
	collection, err := sgf.ParseSgf("(;FF[4]GM[1](;B[qd];W[ob])(;W[pe]))")
	if err != nil {
		panic(err)
	}

	fmt.Println(collection.Sgf(sgf.DefaultSgfFormat))
	fmt.Println()
	fmt.Println(collection.Sgf(sgf.NoNewLinesSgfFormat))
	// Output:
	// (;FF[4]GM[1]
	//     (;B[qd]
	//      ;W[ob])
	//     (;W[pe]))
	//
	// (;FF[4]GM[1](;B[qd];W[ob])(;W[pe]))
}

/*
func ExampleCollection_SwapGameTrees() {
	c, gt1, n1 := sgf.NewCollection()
	_ = n1.NewProperty("FF", "4")

	gt2, n2 := c.NewGameTree()
	_ = n2.NewProperty("FF", "3")

	c.SwapGameTrees(gt1, gt2)

	fmt.Println(c.Sgf(sgf.NoNewLinesSgfFormat))

	// Output: (;FF[3])(;FF[4])
}

func ExampleCollection_SwapGameTreesAt() {
	c, _, n1 := sgf.NewCollection()
	_ = n1.NewProperty("FF", "4")

	_, n2 := c.NewGameTree()
	_ = n2.NewProperty("FF", "3")

	c.SwapGameTreesAt(0, 1)

	fmt.Println(c.Sgf(sgf.NoNewLinesSgfFormat))

	// Output: (;FF[3])(;FF[4])
}

func ExampleCollection_RemoveGameTreeAt() {
	collection, _, node1 := sgf.NewCollection()
	_ = node1.NewProperty("FF", "4")

	_, node2 := c.NewGameTree()
	_ = node2.NewProperty("FF", "3")

	collection.RemoveGameTreeAt(0)

	fmt.Println(c.Sgf(sgf.NoNewLinesSgfFormat))

	// Output: (;FF[3])
}
*/

func ExampleNewCollection_simpleGame() {
	collection, gameTree, rootNode := sgf.NewCollection()

	// Setup root node
	rootNode.NewProperty("FF", "4")
	rootNode.NewProperty("GM", "1")
	rootNode.NewProperty("SZ", "9")

	// Add moves
	gameTree.NewNode().NewProperty("B", "df")
	gameTree.NewNode().NewProperty("W", "cf")
	gameTree.NewNode().NewProperty("B", "ge")
	gameTree.NewNode().NewProperty("W", "ee")

	lastNode := gameTree.NewNode()
	lastNode.NewProperty("W", "bd")
	lastNode.NewProperty("C", "B+R")

	fmt.Println(collection.Sgf(sgf.NoNewLinesSgfFormat))
	// Output:
	// (;FF[4]GM[1]SZ[9];B[df];W[cf];B[ge];W[ee];W[bd]C[B+R])
}

func ExampleNewCollection_goProblem() {
	collection, gameTree, rootNode := sgf.NewCollection()

	// Setup root node
	rootNode.NewProperty("FF", "4")
	rootNode.NewProperty("GM", "1")
	rootNode.NewProperty("SZ", "9")

	// Setup problem
	setupNode := gameTree.NewNode()
	setupNode.NewProperty("AW", "ca", "ea", "bb", "eb", "bc", "bd", "be", "bf")
	setupNode.NewProperty("AB", "fa", "fb", "cc", "dc", "ec", "fc")

	// Add moves
	gameTree.NewNode().NewProperty("B", "cb")
	gameTree.NewNode().NewProperty("W", "ba")
	gameTree.NewNode().NewProperty("B", "da")
	gameTree.NewNode().NewProperty("W", "db")

	lastNode := gameTree.NewNode()
	lastNode.NewProperty("B", "da")
	lastNode.NewProperty("C", "snapback!")

	fmt.Println(collection.Sgf(sgf.NoNewLinesSgfFormat))
	// Output:
	// (;FF[4]GM[1]SZ[9];AW[ca][ea][bb][eb][bc][bd][be][bf]AB[fa][fb][cc][dc][ec][fc];B[cb];W[ba];B[da];W[db];B[da]C[snapback!])
}
