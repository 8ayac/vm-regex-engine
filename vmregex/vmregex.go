// Package vmregex provides a VM regex engine(VM engine).
package vmregex

import (
	"github.com/8ayac/vm-regex-engine/parser"
	"github.com/8ayac/vm-regex-engine/vm"
	"github.com/8ayac/vm-regex-engine/vm/instruction"
	"github.com/8ayac/vm-regex-engine/vm/opcode"
)

// Regexp has a VM and regexp string.
type Regexp struct {
	regexp  string
	runtime *vm.VM
}

// NewRegexp return a new Regexp.
func NewRegexp(re string) *Regexp {
	psr := parser.NewParser(re)
	ast := psr.GetAST()
	bc := ast.Compile()
	bc.AddInst(instruction.NewInst(opcode.Match, 0, nil, nil), bc.N)
	bc.RemoveNOP()

	return &Regexp{
		regexp:  re,
		runtime: vm.NewVM(bc),
	}
}

// Compile is a wrapper function of NewRegexp().
func Compile(re string) *Regexp {
	return NewRegexp(re)
}

// Match returns whether the input string matches the regular expression.
func (re *Regexp) Match(s string) bool {
	return re.runtime.Run(s + "\x00")
}
