// Package parser implements function to parse the regular expressions.
package parser

import (
	"fmt"
	"log"

	"github.com/8ayac/vm-regex-engine/lexer"
	"github.com/8ayac/vm-regex-engine/node"
	"github.com/8ayac/vm-regex-engine/token"
)

// Parser has a slice of tokens to parse, and now looking token.
type Parser struct {
	tokens []*token.Token
	look   *token.Token
}

// NewParser returns a new Parser with the tokens to
// parse that were obtained by scanning.
func NewParser(s string) *Parser {
	p := &Parser{
		tokens: lexer.NewLexer(s).Scan(),
	}
	p.move()
	return p
}

// GetAST returns the root node of AST obtained by parsing.
func (psr *Parser) GetAST() node.Node {
	ast := psr.expression()
	return ast
}

// move updates the now looking token to the next token in token slice.
// If token slice is empty, will set token.EOF as now looking token.
func (psr *Parser) move() {
	if len(psr.tokens) == 0 {
		psr.look = token.NewToken('\x00', token.EOF)
	} else {
		psr.look = psr.tokens[0]
		psr.tokens = psr.tokens[1:]
	}
}

// moveWithValidation execute move() with validating whether
// now looking Token type is an expected (or not).
func (psr *Parser) moveWithValidation(expect token.Type) {
	if psr.look.Ty != expect {
		err := fmt.Sprintf("[syntax error] expect:\x1b[31m%s\x1b[0m actual:\x1b[31m%s\x1b[0m", expect, psr.look.Ty)
		log.Fatal(err)
	}
	psr.move()
}

// expression -> subexpr
func (psr *Parser) expression() node.Node {
	nd := psr.subexpr()
	psr.moveWithValidation(token.EOF)
	return nd
}

// subexpr -> subexpr '|' seq | seq
// (
//	subexpr  -> seq _subexpr
//	_subexpr -> '|' seq _subexpr | ε
// )
func (psr *Parser) subexpr() node.Node {
	nd := psr.seq()
	for {
		if psr.look.Ty == token.UNION {
			psr.moveWithValidation(token.UNION)
			nd2 := psr.seq()
			nd = node.NewUnion(nd, nd2)
		} else {
			break
		}
	}
	return nd
}

// seq -> subseq | ε
func (psr *Parser) seq() node.Node {
	if psr.look.Ty == token.LPAREN || psr.look.Ty == token.CHARACTER || psr.look.Ty == token.ANY {
		return psr.subseq()
	}
	return node.NewEpsilon()
}

// subseq -> subseq sufope | sufope
// (
//	subseq  -> sufope _subseq
//	_subseq -> sufope _subseq | ε
// )
func (psr *Parser) subseq() node.Node {
	nd := psr.sufope()
	if psr.look.Ty == token.LPAREN || psr.look.Ty == token.CHARACTER || psr.look.Ty == token.ANY {
		nd2 := psr.subseq()
		return node.NewConcat(nd, nd2)
	}
	return nd
}

// sufope -> factor ('*'|'+'|'?') | factor
func (psr *Parser) sufope() node.Node {
	nd := psr.factor()
	switch psr.look.Ty {
	case token.STAR:
		psr.move()
		return node.NewStar(nd)
	case token.PLUS:
		psr.move()
		return node.NewPlus(nd)
	case token.QUESTION:
		psr.move()
		return node.NewQuestion(nd)
	}
	return nd
}

// factor -> '(' subexpr ')' | ANY | CHARACTER |
func (psr *Parser) factor() node.Node {
	switch psr.look.Ty {
	case token.LPAREN:
		psr.moveWithValidation(token.LPAREN)
		nd := psr.subexpr()
		psr.moveWithValidation(token.RPAREN)
		return nd
	case token.ANY:
		nd := node.NewAny()
		psr.moveWithValidation(token.ANY)
		return nd
	default:
		nd := node.NewCharacter(psr.look.V)
		psr.moveWithValidation(token.CHARACTER)
		return nd
	}
}
