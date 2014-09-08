/*
Package sgf provides Smart Game Format (SGF) file parsing implementation.

Example file parsing:
	collection, err := sgf.ParseSgfFile("kogo.sgf")
	if err != nil {
		panic(err)
	}

	// Iterate all game trees
	for _, gameTree := range collection.GameTrees {
	       // Iterate all game tree's nodes
	       for _, node := range gameTree.Nodes {
		       // Iterate all node's properties
		       for _, property := range node.Properties {
			       // ...
		       }
	       }

		// Iterate all child game trees
		for _, childGameTree := range gameTree.GameTrees {
			// ...
		}
	}

Example collection creation:

	collection, gameTree, node := sgf.NewCollection()

	// Add some properties to root node
	node.NewProperty("FF", "4")
	node.NewProperty("GM", "1")

	// Create a new node and add some properties to it
	gameTree.NewNode().NewProperty("AW", "bb", "cb", "cc")

	// Add a child game tree
	childGameTree, node2 := gameTree.NewGameTree()
	node2.NewProperty("B", "af")

	// Add new GameTree to the Collection
	newGameTree, node3 := collection.NewGameTree()

	// ...

Converting collection to SGF:
	collection.Sgf(sgf.DefaultSgfFormat)

Collection manipulation:
	gt1 := collection.NewGameTree()
	gt2 := collection.NewGameTree()
	gt3 := collection.NewGameTree()

	// Swap positions of gt1 and gt3. Current order: [gt3, gt2, gt1].
	collection.SwapGameTree(gt1, gt3)

	// Swap positions of gt2 and gt3. Current order: [gt2, gt3, gt1].
	collection.SwapGameTreeAt(0, 1)

	// Remove gt2 from the collection. Current collections: [gt3, gt1].
	collection.RemoveGameTree(gt2)

	// Remove gt1 from the collection by index. Current collections: [gt3].
	collection.RemoveGameTreeAt(1)
*/
package sgf
