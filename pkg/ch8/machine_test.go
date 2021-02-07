package ch8

import (
	"testing"
)

//===========================================================================//
// Decode Tests
//===========================================================================//

func TestDecode(t *testing.T) {
	t.Run("TestDecode", func(*testing.T) {
		t.Run("getX", testGetX)
		t.Run("getY", testGetY)
		t.Run("getN", testGetN)
		t.Run("getOp", testGetOp)
		t.Run("getKK", testGetKK)
		t.Run("getNNN", testGetNNN)
	})
}

func testGetX(t *testing.T) {
	if getX(0x631f) != 0x3 {
		t.Fail()
	}
	if getX(0x9ac7) != 0xa {
		t.Fail()
	}
}

func testGetY(t *testing.T) {
	if getY(0x631f) != 0x1 {
		t.Fail()
	}
	if getY(0x9ac7) != 0xc {
		t.Fail()
	}
}

func testGetN(t *testing.T) {
	if getN(0x631f) != 0xf {
		t.Fail()
	}
	if getN(0x9ac7) != 0x7 {
		t.Fail()
	}
}

func testGetOp(t *testing.T) {
	if getOp(0x631f) != 0x6 {
		t.Fail()
	}
	if getOp(0x9ac7) != 0x9 {
		t.Fail()
	}
}

func testGetKK(t *testing.T) {
	if getKK(0x631f) != 0x1f {
		t.Fail()
	}
	if getKK(0x9ac7) != 0xc7 {
		t.Fail()
	}
}

func testGetNNN(t *testing.T) {
	if getNNN(0x631f) != 0x31f {
		t.Fail()
	}
	if getNNN(0x9ac7) != 0xac7 {
		t.Fail()
	}
}
