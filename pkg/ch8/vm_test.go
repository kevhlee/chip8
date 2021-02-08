package ch8

import (
	"testing"
)

//===========================================================================//
// Decode Tests
//===========================================================================//

var vm *VirtualMachine = NewVirtualMachine()

func TestDecode(t *testing.T) {
	t.Run("TestDecode", func(*testing.T) {
		t.Run("decodeX", testDecodeX)
		t.Run("decodeY", testDecodeY)
		t.Run("decodeN", testDecodeN)
		t.Run("decodeOp", TestDecodeOp)
		t.Run("decodeKK", testDecodeKK)
		t.Run("decodeNNN", testDecodeNNN)
	})
}

func testDecodeX(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeX() != 0x3 {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeX() != 0xa {
		t.Fail()
	}
}

func testDecodeY(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeY() != 0x1 {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeY() != 0xc {
		t.Fail()
	}
}

func testDecodeN(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeN() != 0xf {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeN() != 0x7 {
		t.Fail()
	}
}

func TestDecodeOp(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeOp() != 0x6 {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeOp() != 0x9 {
		t.Fail()
	}
}

func testDecodeKK(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeKK() != 0x1f {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeKK() != 0xc7 {
		t.Fail()
	}
}

func testDecodeNNN(t *testing.T) {
	if vm.Opcode = 0x631f; vm.decodeNNN() != 0x31f {
		t.Fail()
	}
	if vm.Opcode = 0x9ac7; vm.decodeNNN() != 0xac7 {
		t.Fail()
	}
}
