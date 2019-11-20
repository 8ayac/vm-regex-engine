// Package bytecode provides the bytecode structure and some utilities for convenient use of bytecodes.
package bytecode

import (
	"fmt"
	"github.com/8ayac/vm-regex-engine/vm/instruction"
	"github.com/8ayac/vm-regex-engine/vm/opcode"
)

// BC represents a line of instructions.
type BC struct {
	N    int // size
	Code []*instruction.Inst
}

func (bc BC) String() string {
	s := ""
	for i, inst := range bc.Code {
		s += fmt.Sprintf("|%02d:(%p)| %v\n", i, inst, inst)
	}
	return s
}

// NewByteCode returns a new BC.
func NewByteCode() *BC {
	return &BC{
		N:    0,
		Code: make([]*instruction.Inst, 0),
	}
}

// IndexOf returns the position of argument instruction in BC.
// If can't find, then return -1.
func (bc BC) IndexOf(sinst *instruction.Inst) int {
	for i, inst := range bc.Code {
		if inst == sinst {
			return i
		}
	}
	return -1
}

// AddCode inserts the argument bytecode into BC.
// The argument integer i is the index where will insert into.
func (bc *BC) AddCode(bc2 BC, i int) {
	if i > bc.N {
		panic("Too large i.")
	}
	bc.Code = append(bc.Code[:i], append(bc2.Code, bc.Code[i:]...)...)
	bc.N += bc2.N
}

// PushCode inserts the argument bytecode into the top of BC.
func (bc *BC) PushCode(bc2 BC) {
	bc.AddCode(bc2, 0)
}

// AddInst adds the argument instruction to BC.
// The argument integer i is the index where will insert to.
func (bc *BC) AddInst(inst *instruction.Inst, i int) {
	if i < bc.N {
		bc.Code = append(bc.Code[:i+1], bc.Code[i:]...)
		bc.Code[i] = inst
	} else {
		bc.Code = append(bc.Code, inst)
	}
	bc.N++
}

// PushInst adds the argument instruction to the top of BC.
func (bc *BC) PushInst(inst *instruction.Inst) {
	bc.AddInst(inst, 0)
}

// Optimize optimizes (or minimize) the bytecode.
func (bc *BC) Optimize() {
	bc.removeNOP()
	bc.optimizeChainedJmp()
	bc.removeUnreachableInst()
}

// removeNOP removes NOP instructions from BC to minimize it.
func (bc *BC) removeNOP() {
	newBC := NewByteCode()
	for _, inst := range bc.Code {
		switch inst.Opcode {
		case opcode.Jmp:
			if inst.X.Opcode == opcode.NOP {
				inst.X = bc.canReachWithOnlyNOP(bc.IndexOf(inst.X))
			}
		case opcode.Split:
			if inst.X.Opcode == opcode.NOP {
				inst.X = bc.canReachWithOnlyNOP(bc.IndexOf(inst.X))
			}
			if inst.Y.Opcode == opcode.NOP {
				inst.Y = bc.canReachWithOnlyNOP(bc.IndexOf(inst.Y))
			}
		}
		if inst.Opcode != opcode.NOP {
			newBC.AddInst(inst, newBC.N)
		}
	}
	bc.Code = newBC.Code
	bc.N = newBC.N
}

func (bc *BC) canReachWithOnlyNOP(from int) *instruction.Inst {
	prog := bc.Code
	for i := from; i < bc.N; i++ {
		if prog[i].Opcode != opcode.NOP {
			return prog[i]
		}
	}
	return nil
}

// optimizeChainedJmp optimizes chained Jmp instruction. (e.g. [00:Jmp 01 -> 01:Jmp 02 -> 02:Jmp 03 -> 03:Jmp Match] to [Jmp 03 -> 03:Jmp Match])
func (bc *BC) optimizeChainedJmp() {
	for _, inst := range bc.Code {
		switch inst.Opcode {
		case opcode.Jmp:
			if inst.X.Opcode == opcode.Jmp {
				inst.X = bc.chainedJmpDst(inst.X)
			}
		case opcode.Split:
			if inst.X.Opcode == opcode.Jmp {
				inst.X = bc.chainedJmpDst(inst.X)
			}
			if inst.Y.Opcode == opcode.Jmp {
				inst.Y = bc.chainedJmpDst(inst.Y)
			}
		}
	}
}

// chainedJmpDst returns the instruction which in the end will be reached from the argument instruction.
func (bc *BC) chainedJmpDst(jmp *instruction.Inst) *instruction.Inst {
	// Use only for optimizing jmp destination.
	if jmp.Opcode != opcode.Jmp {
		panic("the inst opcode is NOT Jmp")
	}

	// Check recursively.
	dst := jmp.X
	for {
		if dst.Opcode != opcode.Jmp {
			break
		}
		dst = dst.X
	}
	return dst
}

// removeUnreachableInst removes all the unreachable instructions from the bytecode.
func (bc *BC) removeUnreachableInst() {
	bc.removeUnreachableJmp()
}

// removeUnreachableJmp removes unreachable Jmp instruction from the bytecode.
func (bc *BC) removeUnreachableJmp() {
	for i := bc.N - 1; i > 0; i-- {
		if bc.Code[i].Opcode == opcode.Jmp && bc.Code[i-1].Opcode == opcode.Jmp && bc.heldAsOperandCnt(bc.Code[i]) == 0 {
			bc.Code = append(bc.Code[:i], bc.Code[i+1:]...)
			bc.N--
		}
	}
}

// heldAsOperandCnt returns how many instructions in bytecode have inst as operand.
func (bc *BC) heldAsOperandCnt(inst *instruction.Inst) int {
	cnt := 0
	for _, inst := range bc.Code {
		if inst.X == inst {
			cnt++
		}
		if inst.Y == inst {
			cnt++
		}
	}
	return cnt
}
