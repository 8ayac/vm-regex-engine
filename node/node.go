// Package node implements some AST nodes.
package node

import (
	"fmt"
	"github.com/8ayac/vm-regex-engine/bytecode"
	"github.com/8ayac/vm-regex-engine/vm/instruction"
	"github.com/8ayac/vm-regex-engine/vm/opcode"
)

// String to identify the type of Node.
const (
	TypeCharacter = "Character"
	TypeUnion     = "Union"
	TypeConcat    = "Concat"
	TypeStar      = "Star"
	TypePlus      = "Plus"
	TypeQuestion  = "Question"
	TypeEpsilon   = "Epsilon" // Empty character
)

// Node is the interface Node implements.
type Node interface {
	// SubtreeString returns a string to which converts
	// a subtree with Node at the top.
	SubtreeString() string

	// Compile returns a byte code fragment for VM.
	Compile() *bytecode.BC
}

// Character represents the Character node.
type Character struct {
	Ty string
	V  rune
}

func (c *Character) String() string {
	return c.SubtreeString()
}

// NewCharacter returns a new Character node.
func NewCharacter(r rune) *Character {
	return &Character{
		Ty: TypeCharacter,
		V:  r,
	}
}

// Compile returns a BC compiled from Character node which VM can execute.
// The BC compiled from an expression 'a' will be like below:
//
// |0| Char 'a'
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (c *Character) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()
	bc.PushInst(instruction.NewInst(opcode.Char, c.V, nil, nil))
	return bc
}

// SubtreeString returns a string to which converts
// a subtree with the Character node at the top.
func (c *Character) SubtreeString() string {
	return fmt.Sprintf("\x1b[32m%s('%s')\x1b[0m", c.Ty, string(c.V))
}

// Union represents the Union node.
type Union struct {
	Ty   string
	Ope1 Node
	Ope2 Node
}

func (u *Union) String() string {
	return u.SubtreeString()
}

// NewUnion returns a new Union node.
func NewUnion(ope1, ope2 Node) *Union {
	return &Union{
		Ty:   TypeUnion,
		Ope1: ope1,
		Ope2: ope2,
	}
}

// Compile returns a BC compiled from Union node which VM can execute.
// The BC compiled from an expression 'a|b' will be like below:
//
// |0| Split 1, 3
// |1| Char 'a'
// |2| Jmp 4
// |3| Char 'b'
// |4| <nop>
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (u *Union) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()

	e1 := u.Ope1.Compile()
	e2 := u.Ope2.Compile()

	l1 := e1.Code[0]
	l2 := e2.Code[0]
	l3 := instruction.NewInst(opcode.NOP, 0, nil, nil)

	bc.PushInst(l3)
	bc.PushCode(*e2)
	bc.PushInst(instruction.NewInst(opcode.Jmp, 0, l3, nil))
	bc.PushCode(*e1)
	bc.PushInst(instruction.NewInst(opcode.Split, 0, l1, l2))

	return bc
}

// SubtreeString returns a string to which converts
// a subtree with the Union node at the top.
func (u *Union) SubtreeString() string {
	return fmt.Sprintf("\x1b[36m%s(\x1b[0m%s\x1b[36m, \x1b[0m%s\x1b[36m)\x1b[0m", u.Ty, u.Ope1.SubtreeString(), u.Ope2.SubtreeString())
}

// Concat represents the Concat node.
type Concat struct {
	Ty   string
	Ope1 Node
	Ope2 Node
}

// Compile returns a BC compiled from Concat node which VM can execute.
// The BC compiled from an expression 'abc' will be like below:
//
// |0| Char 'a'
// |1| Char 'b'
// |2| Char 'c'
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (c *Concat) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()

	e1 := c.Ope1.Compile()
	e2 := c.Ope2.Compile()

	bc.PushCode(*e2)
	bc.PushCode(*e1)

	return bc
}

func (c *Concat) String() string {
	return c.SubtreeString()
}

// NewConcat returns a new Concat node.
func NewConcat(ope1, ope2 Node) *Concat {
	return &Concat{
		Ty:   TypeConcat,
		Ope1: ope1,
		Ope2: ope2,
	}
}

// SubtreeString returns a string to which converts
// a subtree with the Concat node at the top.
func (c *Concat) SubtreeString() string {
	return fmt.Sprintf("\x1b[31m%s(\x1b[0m%s\x1b[31m, \x1b[0m%s\x1b[31m)\x1b[0m", c.Ty, c.Ope1.SubtreeString(), c.Ope2.SubtreeString())
}

// Star represents the Star node.
type Star struct {
	Ty  string
	Ope Node
}

// Compile returns a BC compiled from Star node which VM can execute.
// The BC compiled from an expression 'abc' will be like below:
//
// |00| Split 1, 3
// |01| Char 'a'
// |02| Jmp 0
// |03| <nop>
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (s *Star) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()

	e := s.Ope.Compile()

	l3 := instruction.NewInst(opcode.NOP, 0, nil, nil)
	l2 := e.Code[0]
	l1 := instruction.NewInst(opcode.Split, 0, l2, l3)

	bc.PushInst(l3)
	bc.PushInst(instruction.NewInst(opcode.Jmp, 0, l1, nil))
	bc.PushCode(*e)
	bc.PushInst(l1)

	return bc
}

func (s *Star) String() string {
	return s.SubtreeString()
}

// NewStar returns a new Star node.
func NewStar(ope Node) *Star {
	return &Star{
		Ty:  TypeStar,
		Ope: ope,
	}
}

// SubtreeString returns a string to which converts
// a subtree with the Star node at the top.
func (s *Star) SubtreeString() string {
	return fmt.Sprintf("\x1b[33m%s(%s\x1b[33m)\x1b[0m", s.Ty, s.Ope.SubtreeString())
}

// Plus represents the Plus node.
type Plus struct {
	Ty  string
	Ope Node
}

// Compile returns a BC compiled from Plus node which VM can execute.
// The BC compiled from an expression 'a' will be like below:
//
// |00| Char 'a'
// |01| Split 0, 2
// |02| <nop>
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (p *Plus) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()

	e := p.Ope.Compile()
	l1 := e.Code[0]
	l2 := instruction.NewInst(opcode.NOP, 0, nil, nil)

	bc.PushInst(l2)
	bc.PushInst(instruction.NewInst(opcode.Split, 0, l1, l2))
	bc.PushCode(*e)

	return bc
}

func (p *Plus) String() string {
	return p.SubtreeString()
}

// NewPlus returns a new Plus node.
func NewPlus(ope Node) *Plus {
	return &Plus{
		Ty:  TypePlus,
		Ope: ope,
	}
}

// SubtreeString returns a string to which converts
// a subtree with the Star node at the top.
func (p *Plus) SubtreeString() string {
	return fmt.Sprintf("\x1b[33m%s(%s\x1b[33m)\x1b[0m", p.Ty, p.Ope.SubtreeString())
}

// Question represents the Question node.
type Question struct {
	Ty  string
	Ope Node
}

// Compile returns a BC compiled from Question node which VM can execute.
// The BC compiled from an expression 'a' will be like below:
//
// |00| Split 1, 2
// |01| Char 'a'
// |02| <nop>
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (q *Question) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()

	e := q.Ope.Compile()
	l1 := e.Code[0]
	l2 := instruction.NewInst(opcode.NOP, 0, nil, nil)

	bc.PushInst(l2)
	bc.PushInst(l1)
	bc.PushInst(instruction.NewInst(opcode.Split, 0, l1, l2))

	return bc
}

func (q *Question) String() string {
	return q.SubtreeString()
}

// NewQuestion returns a new Question node.
func NewQuestion(ope Node) *Question {
	return &Question{
		Ty:  TypeQuestion,
		Ope: ope,
	}
}

// SubtreeString returns a string to which converts
// a subtree with the Star node at the top.
func (q *Question) SubtreeString() string {
	return fmt.Sprintf("\x1b[33m%s(%s\x1b[33m)\x1b[0m", q.Ty, q.Ope.SubtreeString())
}

// Epsilon represents the Star node.
type Epsilon struct {
	Ty string
}

// Compile returns a BC compiled from Epsilon node which VM can execute.
// The BC compiled from an expression 'abc' will be like below:
//
// |00| <nop>
//
// Note:
// The bytecode is just a fragment, so when finally give VM it,
// you need to add the instruction of Match to the last of BC.
func (*Epsilon) Compile() *bytecode.BC {
	bc := bytecode.NewByteCode()
	bc.PushInst(instruction.NewInst(opcode.NOP, 0, nil, nil))
	return bc
}

func (e *Epsilon) String() string {
	return e.SubtreeString()
}

// NewEpsilon returns a new Star node.
func NewEpsilon() *Epsilon {
	return &Epsilon{
		Ty: TypeEpsilon,
	}
}

// SubtreeString returns a string to which converts
// a subtree with the Epsilon node at the top.
func (e *Epsilon) SubtreeString() string {
	return fmt.Sprintf("\x1b[37m%s(Empty)\x1b[0m", e.Ty)
}
