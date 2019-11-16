// Package instruction provides the instruction structure and some its methods.
// Using one or more instructions build a bytecode for VM.
package instruction

import (
	"fmt"
	"github.com/8ayac/vm-regex-engine/vm/opcode"
)

// Inst represents a instruction which is executable with the VM.
type Inst struct {
	Opcode opcode.Opcode
	C      rune  // operand for Char
	X      *Inst // operand for Jmp, Split
	Y      *Inst // operand for Split
}

func (inst Inst) String() string {
	switch inst.Opcode {
	case opcode.Char:
		return fmt.Sprintf("Char '%c'", inst.C)
	case opcode.Match:
		return fmt.Sprintf("Match")
	case opcode.Jmp:
		return fmt.Sprintf("Jmp %p(%+v)", inst.X, inst.X)
	case opcode.Split:
		return fmt.Sprintf("Split %p(%+v), %p(%+v)", inst.X, inst.X, inst.Y, inst.Y)
	case opcode.Dummy:
		return fmt.Sprintf("<nop>")
	}
	return ""
}

// NewInst returns a new Inst.
func NewInst(op opcode.Opcode, c rune, x, y *Inst) *Inst {
	return &Inst{
		Opcode: op,
		C:      c,
		X:      x,
		Y:      y,
	}
}
