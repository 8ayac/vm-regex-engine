// Package vm provides a Virtual Machine to execute regular expression matching.
// For Details: https://swtch.com/~rsc/regexp/regexp2.html
package vm

import (
	"github.com/8ayac/vm-regex-engine/bytecode"
	"github.com/8ayac/vm-regex-engine/vm/instruction"
	"github.com/8ayac/vm-regex-engine/vm/opcode"
)

// VM represents a Virtual Machine which executes regular expression matching.
// The VM has a bytecode and one or more threads.
type VM struct {
	bc      bytecode.BC
	threads []*Thread
}

// NewVM returns a new VM for executing argument bytecode.
func NewVM(bc *bytecode.BC) *VM {
	bc.AddInst(instruction.NewInst(opcode.Match, 0, nil, nil), bc.N)

	return &VM{
		bc:      *bc,
		threads: []*Thread{},
	}
}

// AddThread adds a new Thread to the stack of threads in VM.
func (v *VM) AddThread(t *Thread) {
	v.threads = append(v.threads, t)
}

// Run starts to execute regular expression matching for the input string.
// The return value is determined by whether matching is success or not.
// If success then will be true, but not success will be false.
func (v *VM) Run(input string) bool {
	const MAXTHREAD = 10000 //

	prog := v.bc.Code
	ready := [MAXTHREAD]*Thread{}
	nready := 0

	var pc int
	var sp int

	ready[0] = NewThread(pc, sp)
	nready = 1

	for nready > 0 {
		nready--

		pc = ready[nready].PC
		sp = ready[nready].SP

		for {
			switch prog[pc].Opcode {
			case opcode.Char:
				if []rune(input)[sp] != prog[pc].C {
					goto Dead
				}
				sp++
				pc++
				continue
			case opcode.Match:
				return true
			case opcode.Jmp:
				pc = v.bc.IndexOf(prog[pc].X)
				continue
			case opcode.Split:
				if nready >= MAXTHREAD {
					panic("The thread in VM overflowed!")
				}
				ready[nready] = NewThread(v.bc.IndexOf(prog[pc].Y), sp)
				nready++
				pc = v.bc.IndexOf(prog[pc].X)
				continue
			case opcode.NOP:
				pc++
				continue
			}
		}
	Dead:
	}
	return false
}

// Thread represents a thread which has two pointers(program counter/string pointer).
// A program counter (PC) is a register has the information where a instruction which being executed by VM.
// A string pointer (SP) is a register has the information where a character that the VM is looking at.
type Thread struct {
	PC int
	SP int
}

// NewThread returns a new Thread which has the PC and SP set to the value specified by the argument.
func NewThread(pc int, sp int) *Thread {
	return &Thread{
		PC: pc,
		SP: sp,
	}
}
