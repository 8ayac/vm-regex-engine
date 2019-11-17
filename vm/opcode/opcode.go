// Package opcode provides the type Opcode to identify the type of operation code.
package opcode

// Opcode is integer to identify the type of operation code.
type Opcode int

func (op Opcode) String() string {
	switch op {
	case Char:
		return "Char"
	case Match:
		return "Match"
	case Jmp:
		return "Jmp"
	case Split:
		return "Split"
	case NOP:
		return "NOP"
	}
	return ""
}

const (
	Char Opcode = iota
	Match
	Jmp
	Split
	NOP
)
