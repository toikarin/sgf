package sgf

import (
	"errors"
	"fmt"
)

// Token Type
type tokenType int

const (
	tokenTypeGameTreeStart tokenType = iota
	tokenTypeGameTreeEnd
	tokenTypeNode
	tokenTypePropertyIdent
	tokenTypePropertyValue
)

// Lexer states
const (
	lexerStateOnlyControl = iota
	lexerStatePropertyValue
)

// Parser states
const (
	parserStateCollection = iota
	parserStateGameTree
	parserStateValue
	parserStateNode
)

type lexeme struct {
	tokenType tokenType
	data      string
	line      int
	position  int
}

func lexicalAnalysis(data string) ([]lexeme, error) {
	retval := []lexeme{}
	curData := ""
	escapedText := false
	line := 0
	position := 0

	var lexerState = lexerStateOnlyControl

	for _, c := range data {
		switch lexerState {
		case lexerStateOnlyControl:
			{
				position++

				// control chars
				switch c {
				case ' ', '\t', '\v', '\r', '\n':
					if curData != "" {
						return nil, createLexerError(fmt.Sprintf("Invalid character %c", c), line, position)
					}

					if c == '\n' || c == '\r' {
						line++
					}
				case '(', ')', ';':
					if curData != "" {
						retval = append(retval, lexeme{tokenTypePropertyIdent, curData, line, position})
					}

					retval = append(retval, lexeme{runeToTokenType(c), "", line, position})
					curData = ""
				case '[':
					if curData != "" {
						retval = append(retval, lexeme{runeToTokenType(c), curData, line, position})
					}
					curData = ""
					lexerState = lexerStatePropertyValue
				default:
					if c < 'A' || c > 'Z' {
						return nil, createLexerError(fmt.Sprintf("Invalid character %c", c), line, position)
					}

					curData += string(c)
				}
			}
		case lexerStatePropertyValue:
			{
				position++
				if c == '\\' {
					if escapedText {
						curData += string(c)
						escapedText = false
					} else {
						escapedText = true
					}
				} else if c == ']' && !escapedText {
					retval = append(retval, lexeme{tokenTypePropertyValue, curData, line, position})

					lexerState = lexerStateOnlyControl
					curData = ""
					escapedText = false
				} else {
					curData += string(c)
					escapedText = false
				}
			}
		}
	}

	if curData != "" {
		if lexerState == lexerStateOnlyControl {
			retval = append(retval, lexeme{tokenTypePropertyIdent, curData, line, position})
		} else {
			return nil, createLexerError("value left open", line, position)
		}
	}

	return retval, nil
}

func runeToTokenType(r rune) tokenType {
	switch r {
	case '(':
		return tokenTypeGameTreeStart
	case ')':
		return tokenTypeGameTreeEnd
	case ';':
		return tokenTypeNode
	case '[':
		return tokenTypePropertyIdent
	default:
		panic("Invalid rune: " + string(r))
	}
}

func parse(lexemes []lexeme) (*Collection, error) {
	c := Collection{}
	gameTreeStack := make([]*GameTree, 0)
	var curGameTree *GameTree
	var curNode *Node
	var curProperty *Property

	state := parserStateCollection

	for _, l := range lexemes {
		switch state {
		// This state occurs only on the root of the SGF files. Can happen multiple times in a single file but must
		// always start a new game tree.
		// Example file containing two lines has two different game trees in one collection:
		//   (;FF[4]...)
		//   (;FF[4]...)
		//   ^
		case parserStateCollection:
			if l.tokenType != tokenTypeGameTreeStart {
				return nil, createParserError("Collection must start with a new game tree.", l)
			}

			// Create a new game tree
			curGameTree = &GameTree{}
			// Add game tree to collection
			c.GameTrees = append(c.GameTrees, curGameTree)
			// Change state
			state = parserStateGameTree
		// Parsing of a game tree. Game tree always starts always with a new node.
		// Example file:
		//   (;FF[4](;B[])(;B[])...)
		//    ^      ^     ^
		case parserStateGameTree:
			if l.tokenType != tokenTypeNode {
				return nil, createParserError("New node must be next after game tree has started.", l)
			}

			curNode = &Node{}
			curGameTree.Nodes = append(curGameTree.Nodes, curNode)
			state = parserStateNode
		// Parsing of a node.
		// Node can have a zero or more properties. Properties always start with an idend and they have at least one
		// value. After a node we can either start a new node, start a new game tree (part of the current game tree)
		// or close the current game tree.
		// Few examples:
		// (;)
		// (;;)
		// (;AB[bb:ee])
		// (;FF[4];GM[1])
		// (;AW[bb][ee][dc][cd])
		// (;FF[4](;B[1]))
		case parserStateNode:
			switch l.tokenType {
			case tokenTypePropertyIdent:
				// New property starts
				curProperty = &Property{Ident: l.data}
				curNode.Properties = append(curNode.Properties, curProperty)

				// Next must come the value
				state = parserStateValue
			case tokenTypePropertyValue:
				// value can only come after ident
				if curProperty == nil {
					return nil, createParserError("Cannot have property value without property ident.", l)
				}

				// An extra value to current property
				curProperty.Values = append(curProperty.Values, l.data)
			// New node starts after current node
			case tokenTypeNode:
				curNode = &Node{}
				curGameTree.Nodes = append(curGameTree.Nodes, curNode)
				curProperty = nil
			// New game tree starts
			case tokenTypeGameTreeStart:
				// clean up node related state
				curProperty = nil
				curNode = nil

				// create new game tree
				newGameTree := &GameTree{}
				// Append game tree to a current game tree as a child
				curGameTree.GameTrees = append(curGameTree.GameTrees, newGameTree)
				// Add current game tree to stack
				gameTreeStack = append(gameTreeStack, curGameTree)
				// Swap current game tree to new one
				curGameTree = newGameTree

				// Next token must be a node
				state = parserStateGameTree
			// Current game tree ends
			case tokenTypeGameTreeEnd:
				// Clean up node related state
				curProperty = nil
				curNode = nil

				// if stack is empty go to a collection state (whole new game tree must be started)
				if len(gameTreeStack) == 0 {
					curGameTree = nil
					state = parserStateCollection
				} else {
					// take game tree from the stack
					curGameTree = gameTreeStack[len(gameTreeStack)-1]
					gameTreeStack = gameTreeStack[:len(gameTreeStack)-1]
				}
			}
		// Used to make sure there always is value after ident
		case parserStateValue:
			// value can only come after ident
			if l.tokenType != tokenTypePropertyValue {
				return nil, createParserError("After property ident there must be a value.", l)
			}

			// Add value to current property
			curProperty.Values = append(curProperty.Values, l.data)
			state = parserStateNode
		}
	}

	if state != parserStateCollection {
		return nil, createParserError("Game tree did not close properly.", lexemes[len(lexemes)-1])
	}

	return &c, nil
}

func createLexerError(msg string, line, position int) error {
	return errors.New(fmt.Sprintf("%s [line %d, position %d]", msg, line, position))
}

func createParserError(msg string, lexeme lexeme) error {
	return errors.New(fmt.Sprintf("%s [line %d, position %d]", msg, lexeme.line, lexeme.position))
}
