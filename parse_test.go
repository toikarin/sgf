package sgf

import (
	"testing"
)

func ltype(tt tokenType) lexeme {
	return lexeme{tt, "", 0, 0}
}

func lvalue(tt tokenType, value string) lexeme {
	return lexeme{tt, value, 0, 0}
}

func TestLexicalAnalysis(t *testing.T) {
	var okTests = []struct {
		data   string
		wanted []lexeme
	}{
		// empty
		{"", []lexeme{}},
		// minimal example
		{"(;)", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			ltype(tokenTypeGameTreeEnd),
		}},
		// node with simple properties
		{"(;FF[4]GM[1])", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "4"),
			lvalue(tokenTypePropertyIdent, "GM"),
			lvalue(tokenTypePropertyValue, "1"),
			ltype(tokenTypeGameTreeEnd),
		}},
		// few nodes with properties
		{"(;FF[4];GM[1])", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "4"),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "GM"),
			lvalue(tokenTypePropertyValue, "1"),
			ltype(tokenTypeGameTreeEnd),
		}},
		// child game trees
		{"(;FF[4];GM[1](;PB[Black]))", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "4"),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "GM"),
			lvalue(tokenTypePropertyValue, "1"),
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "PB"),
			lvalue(tokenTypePropertyValue, "Black"),
			ltype(tokenTypeGameTreeEnd),
			ltype(tokenTypeGameTreeEnd),
		}},
		// spaces are ignored
		{"\t\r\n( ;  FF[4]\n\nGM[1]\t[2]\r \t;\t(\t;\t)\t)\t", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "4"),
			lvalue(tokenTypePropertyIdent, "GM"),
			lvalue(tokenTypePropertyValue, "1"),
			lvalue(tokenTypePropertyValue, "2"),
			ltype(tokenTypeNode),
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			ltype(tokenTypeGameTreeEnd),
			ltype(tokenTypeGameTreeEnd),
		}},
		// spaces are not ignored in property values
		{"(;FF[\t\n4\r])", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "\t\n4\r"),
			ltype(tokenTypeGameTreeEnd),
		}},
		// escaped text
		{"(;FF[\\]]FF[\\\\])", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "]"),
			lvalue(tokenTypePropertyIdent, "FF"),
			lvalue(tokenTypePropertyValue, "\\"),
			ltype(tokenTypeGameTreeEnd),
		}},
		// illegal, but ok to lex
		{"FF(;)", []lexeme{
			lvalue(tokenTypePropertyIdent, "FF"),
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			ltype(tokenTypeGameTreeEnd),
		}},
		{"(;)FF", []lexeme{
			ltype(tokenTypeGameTreeStart),
			ltype(tokenTypeNode),
			ltype(tokenTypeGameTreeEnd),
			lvalue(tokenTypePropertyIdent, "FF"),
		}},
	}

	for _, test := range okTests {
		lexemes, err := lexicalAnalysis(test.data)
		if err != nil {
			t.Errorf("lexicalAnalysis(%s) returned error.", test.data)
			continue
		}

		for i, wanted := range test.wanted {
			if len(lexemes) < i {
				t.Errorf("lexicalAnalysis(%s) not enough results.", test.data)
				break
			}

			if lexemes[i].tokenType != wanted.tokenType {
				t.Errorf("lexicalAnalysis(%s) token type mismatch at index %d. wanted: %d, got: %d.", test.data, i, wanted.tokenType, lexemes[i].tokenType)
			}
			if lexemes[i].data != wanted.data {
				t.Errorf("lexicalAnalysis(%s) data mismatch at index %d. wanted: %s, got: %s.", test.data, i, wanted.data, lexemes[i].data)
			}
		}

		if len(lexemes) > len(test.wanted) {
			t.Errorf("lexicalAnalysis(%s) returned more lexemes than expected. wanted: %d, got: %d.", test.data, len(test.wanted), len(lexemes))
		}
	}
}

func TestLexicalAnalysisErrors(t *testing.T) {
	var errTests = []string{
		"abc",
		"'",
		"FF[4",
		"FF [4]",
	}

	for _, test := range errTests {
		_, err := lexicalAnalysis(test)
		if err == nil {
			t.Errorf("lexicalAnalysis(%s) did not return error.", test)
		}
	}
}

func TestParse(t *testing.T) {
	var okTests = []struct {
		data               string
		expectedCollection Collection
	}{
		{"(;)",
			Collection{
				[]*GameTree{
					&GameTree{
						[]*Node{
							&Node{},
						},
						[]*GameTree{},
					},
				},
			},
		},
		{"(;;)",
			Collection{
				[]*GameTree{
					&GameTree{
						[]*Node{
							&Node{},
							&Node{},
						},
						[]*GameTree{},
					},
				},
			},
		},
		{"(;FF[4])",
			Collection{
				[]*GameTree{
					&GameTree{
						[]*Node{
							&Node{[]*Property{
								&Property{"FF", []string{"4"}},
							}},
						},
						[]*GameTree{},
					},
				},
			},
		},
		{"(;FF[4][5])",
			Collection{
				[]*GameTree{
					&GameTree{
						[]*Node{
							&Node{[]*Property{
								&Property{"FF", []string{"4", "5"}},
							}},
						},
						[]*GameTree{},
					},
				},
			},
		},
		{"(;FF[4]SZ[1](;PB[Black]))",
			Collection{
				[]*GameTree{
					&GameTree{
						[]*Node{
							&Node{[]*Property{
								&Property{"FF", []string{"4"}},
								&Property{"SZ", []string{"1"}},
							}},
						},
						[]*GameTree{
							&GameTree{
								[]*Node{
									&Node{[]*Property{
										&Property{"PB", []string{"Black"}},
									}},
								},
								[]*GameTree{},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range okTests {
		expectedCollection := test.expectedCollection
		actualCollection, err := ParseSgf(test.data)

		if err != nil {
			t.Errorf("ParseSgf(%s) returned error.", test.data)
			continue
		}

		for i, expectedGameTree := range expectedCollection.GameTrees {
			if len(actualCollection.GameTrees) < i {
				t.Errorf("lexicalAnalysis(%s) foo.", test.data)
				break
			}

			actualGameTree := actualCollection.GameTrees[i]
			equalGameTrees(t, test.data, i, expectedGameTree, actualGameTree)
		}

		if len(actualCollection.GameTrees) > len(expectedCollection.GameTrees) {
			t.Errorf("ParseSGF(%s) returned more game trees than expected. wanted: %d, got: %d.", test.data, len(expectedCollection.GameTrees), len(actualCollection.GameTrees))
		}
	}
}

func equalGameTrees(t *testing.T, testData string, i int, expectedGameTree, actualGameTree *GameTree) {
	// Nodes
	for ii, expectedNode := range expectedGameTree.Nodes {
		if len(actualGameTree.Nodes) < ii {
			t.Errorf("lexicalAnalysis(%s) bar.", testData)
			break
		}

		actualNode := actualGameTree.Nodes[ii]

		// Properties
		for iii, expectedProperty := range expectedNode.Properties {
			if len(actualNode.Properties)-1 < iii {
				t.Errorf("lexicalAnalysis(%s) baz.", testData)
				break
			}

			actualProperty := actualNode.Properties[iii]

			if actualProperty.Ident != expectedProperty.Ident {
				t.Errorf("ParseSgf(%s) property mismatch. Property GameTree[%d].Node[%d].Property[%d]. wanted: %s, got: %s.", testData, i, ii, iii, expectedProperty.Ident, actualProperty.Ident)
			}

			for iiii, expectedValue := range expectedProperty.Values {
				actualValue := actualProperty.Values[iiii]

				if actualValue != expectedValue {
					t.Errorf("ParseSgf(%s) property value mismatch. Property GameTree[%d].Node[%d].Property[%d].Value[%d]. wanted: %s, got: %s.", testData, i, ii, iii, iiii, expectedValue, actualValue)
				}
			}
		}

		if len(actualNode.Properties) > len(expectedNode.Properties) {
			t.Errorf("ParseSGF(%s) returned more properties than expected. wanted: %d, got: %d.", testData, len(expectedNode.Properties), len(actualNode.Properties))
		}
	}

	if len(actualGameTree.Nodes) > len(expectedGameTree.Nodes) {
		t.Errorf("ParseSGF(%s) returned more nodes than expected. wanted: %d, got: %d.", testData, len(expectedGameTree.Nodes), len(actualGameTree.Nodes))
	}

	// GameTrees
	for iii, expectedChildGameTree := range expectedGameTree.GameTrees {
		actualChildGameTree := actualGameTree.GameTrees[iii]

		equalGameTrees(t, testData, i, expectedChildGameTree, actualChildGameTree)
	}

	if len(actualGameTree.GameTrees) > len(expectedGameTree.GameTrees) {
		t.Errorf("ParseSGF(%s) returned more child game trees than expected. wanted: %d, got: %d.", testData, len(expectedGameTree.GameTrees), len(actualGameTree.GameTrees))
	}
}

func TestParseErrors(t *testing.T) {
	var errTests = []string{
		"abc",   // Does not lex
		"F",     // Does not start with new game tree
		"(F",    // No new node after game tree
		"(;[V]", // Property value without ident
		"(;FF;", // Property ident without value
		"(;",    // Missing end of game tree
	}

	for _, test := range errTests {
		_, err := ParseSgf(test)
		if err == nil {
			t.Errorf("ParseSgf(%s) did not return error.", test)
		}
	}
}

func TestRuneToTokenPanics(t *testing.T) {
	panicked := fnPanics(func() { runeToTokenType('a') })

	if !panicked {
		t.Errorf(" runeToTokenType('a') should have panicked")
	}
}
