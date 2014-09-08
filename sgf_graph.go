package sgf

/*
import (
	"fmt"
)

func GraphCollection(collection *Collection) {
	for _, gt := range collection.GameTrees {
		graphGameTree(gt, 0, true)
	}
	fmt.Println()
}

func graphGameTree(gameTree *GameTree, indentation int, firstNode bool) {
	if firstNode {
		for i := 0; i < indentation; i++ {
			if i+2 == indentation {
				fmt.Print("+")
			} else if i+1 == indentation {
				fmt.Print("-")
			} else {
				fmt.Print(" ")
			}
		}
	}

	for _, _ = range gameTree.Nodes {
		if !firstNode {
			fmt.Print("-")
			indentation += 2
		}
		fmt.Print("o")
		firstNode = false
	}
	for i, childGameTree := range gameTree.GameTrees {
		if i == 0 {
			graphGameTree(childGameTree, indentation, firstNode)
		} else {
			fmt.Println()
			graphGameTree(childGameTree, indentation+2, true)
		}
	}
}
*/
