// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Provides support for dealing with EVM assembly instructions (e.g., disassembling them).
package asm

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

// Iterator for disassembled EVM instructions
type instructionIterator struct {
	code    []byte
	pc      uint64
	arg     []byte
	op      vm.OpCode
	error   error
	started bool
}

// Create a new instruction iterator.
func NewInstructionIterator(code []byte) *instructionIterator { log.DebugLog()
	it := new(instructionIterator)
	it.code = code
	return it
}

// Returns true if there is a next instruction and moves on.
func (it *instructionIterator) Next() bool { log.DebugLog()
	if it.error != nil || uint64(len(it.code)) <= it.pc {
		// We previously reached an error or the end.
		return false
	}

	if it.started {
		// Since the iteration has been already started we move to the next instruction.
		if it.arg != nil {
			it.pc += uint64(len(it.arg))
		}
		it.pc++
	} else {
		// We start the iteration from the first instruction.
		it.started = true
	}

	if uint64(len(it.code)) <= it.pc {
		// We reached the end.
		return false
	}

	it.op = vm.OpCode(it.code[it.pc])
	if it.op.IsPush() {
		a := uint64(it.op) - uint64(vm.PUSH1) + 1
		u := it.pc + 1 + a
		if uint64(len(it.code)) <= it.pc || uint64(len(it.code)) < u {
			it.error = fmt.Errorf("incomplete push instruction at %v", it.pc)
			return false
		}
		it.arg = it.code[it.pc+1 : u]
	} else {
		it.arg = nil
	}
	return true
}

// Returns any error that may have been encountered.
func (it *instructionIterator) Error() error { log.DebugLog()
	return it.error
}

// Returns the PC of the current instruction.
func (it *instructionIterator) PC() uint64 { log.DebugLog()
	return it.pc
}

// Returns the opcode of the current instruction.
func (it *instructionIterator) Op() vm.OpCode { log.DebugLog()
	return it.op
}

// Returns the argument of the current instruction.
func (it *instructionIterator) Arg() []byte { log.DebugLog()
	return it.arg
}

// Pretty-print all disassembled EVM instructions to stdout.
func PrintDisassembled(code string) error { log.DebugLog()
	script, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			fmt.Printf("%06v: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			fmt.Printf("%06v: %v\n", it.PC(), it.Op())
		}
	}
	return it.Error()
}

// Return all disassembled EVM instructions in human-readable format.
func Disassemble(script []byte) ([]string, error) { log.DebugLog()
	instrs := make([]string, 0)

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			instrs = append(instrs, fmt.Sprintf("%06v: %v 0x%x\n", it.PC(), it.Op(), it.Arg()))
		} else {
			instrs = append(instrs, fmt.Sprintf("%06v: %v\n", it.PC(), it.Op()))
		}
	}
	if err := it.Error(); err != nil {
		return nil, err
	}
	return instrs, nil
}
