package sgf

// Collection contains one or more GameTrees.
type Collection struct {
	GameTrees []*GameTree // at least one
}

// GameTree contains one or more Nodes and zero or more child GameTrees.
type GameTree struct {
	Nodes     []*Node     // at least one
	GameTrees []*GameTree // zero or more
}

// A node contains zero or more Properties.
type Node struct {
	Properties []*Property // zero or more
}

// Property contains ident string and one or more values.
type Property struct {
	Ident  string
	Values []string // at least one
}
