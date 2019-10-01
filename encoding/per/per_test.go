package per

import (
	"fmt"
	"testing"
)

func compareSlice(actual, expect []uint8) bool {
	if len(actual) != len(expect) {
		return false
	}
	for i := 0; i < len(actual); i++ {
		if actual[i] != expect[i] {
			return false
		}
	}
	fmt.Printf("")
	return true
}

func TestMergeBitField(t *testing.T) {
	in1 := []uint8{}
	in2 := []uint8{0x00, 0x11}
	out, outlen := MergeBitField(nil, 0, in2, 8)
	expect := []uint8{0x11, 0x00}
	expectlen := 8
	if expectlen != outlen || compareSlice(expect, out) == false {
		t.Errorf("bitlen expect: %d, actual %d", expectlen, outlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, out)
	}

	/* XXX need to fix here */
	in1 = []uint8{0x80, 0x80} //b1000 0000 1000 0000
	in2 = []uint8{0x00, 0x11} //b0000 0000 0001 0001
	out, outlen = MergeBitField(in1, 9, in2, 9)
	expect = []uint8{0x80, 0x84, 0x40} //b1000 0000 1000 0100 01xx xxxx
	expectlen = 18
	if expectlen != outlen || compareSlice(expect, out) == false {
		t.Errorf("bitlen expect: %d, actual %d", expectlen, outlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, out)
	}
}

func TestShiftLeft(t *testing.T) {
	in := []uint8{0x00, 0x11, 0x22}
	expect := []uint8{0x01, 0x12, 0x20}
	out := ShiftLeft(in, 4)
	if compareSlice(expect, out) == false {
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, out)
	}
}

func TestShiftRight(t *testing.T) {
	in := []uint8{0x00, 0x11, 0x22}
	expect := []uint8{0x00, 0x01, 0x12}
	out := ShiftRight(in, 4)
	if compareSlice(expect, out) == false {
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, out)
	}
}

// 10.5
func TestEncConstrainedWholeNumber(t *testing.T) {
	v, bitlen, err := EncConstrainedWholeNumber(256, 0, 255)
	if err == nil {
		t.Errorf("EncConstrainedWholeNumber: unexpected error")
	}

	v, bitlen, err = EncConstrainedWholeNumber(1, 0, 0)
	expect := []uint8{}
	if bitlen != 0 {
		t.Errorf("bitlen expect: %d, actual %d", 4, bitlen)
	}

	v, bitlen, err = EncConstrainedWholeNumber(1, 0, 7)
	expect = []uint8{0x01}
	if bitlen != 4 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 4, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncConstrainedWholeNumber(128, 0, 255)
	expect = []uint8{128}
	if bitlen != 8 || compareSlice(expect, v) == false {
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncConstrainedWholeNumber(256, 0, 65535)
	expect = []uint8{1, 0}
	if bitlen != 16 || compareSlice(expect, v) == false {
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, v)
	}
}

// 12
func TestEncInteger(t *testing.T) {
	v, bitlen, err := EncInteger(3, 0, 2, true)
	if err == nil {
		t.Errorf("EncInteger: unexpected error")
	}

	v, bitlen, err = EncInteger(2, 2, 2, false)
	if bitlen != 0 && len(v) == 0 {
		t.Errorf("bitlen expect: %d, actual %d", 2, bitlen)
	}

	v, bitlen, err = EncInteger(2, 2, 2, true)
	expect := []uint8{0x00}
	if bitlen != 1 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 2, bitlen)
	}

	v, bitlen, err = EncInteger(128, 0, 255, false)
	expect = []uint8{128}
	if bitlen != 8 || compareSlice(expect, v) == false {
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncInteger(1, 0, 7, true)
	expect = []uint8{0x01}
	if bitlen != 5 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 4, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncInteger(128, 0, 255, true)
	expect = []uint8{0x00, 128}
	if bitlen != 16 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 16, bitlen)
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncInteger(256, 0, 65535, false)
	expect = []uint8{1, 0}
	if bitlen != 16 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 16, bitlen)
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, v)
	}

}

// 13
func TestEncEnumerated(t *testing.T) {
	v, bitlen, err := EncEnumerated(3, 0, 2, false)
	if err == nil {
		t.Errorf("EncEnumerated: unexpected error")
	}

	v, bitlen, err = EncEnumerated(2, 0, 2, false)
	expect := []uint8{0x02}
	if bitlen != 2 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 2, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	v, bitlen, err = EncEnumerated(1, 0, 2, true)
	expect = []uint8{0x01}
	if bitlen != 3 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 3, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}
}

// 13
func TestEncSequence(t *testing.T) {
	v, bitlen, err := EncSequence(false, 8, 0x00)
	if err == nil {
		t.Errorf("EncSequence: unexpected error")
	}

	v, bitlen, err = EncSequence(true, 1, 0x00)
	expect := []uint8{0x00}
	if bitlen != 2 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 2, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

}

/*
func TestChoice(t *testing.T) {

	expect := []uint8{0x00}
	actual := EncChoice(0)
	if compareSlice(actual, expect) == false {
		t.Errorf("expect: 0x%02x, actual 0x%02x", expect, actual)
	}
}
*/

func TestBitString(t *testing.T) {

	in := []uint8{}
	v, bitlen, err := EncBitString(in, 0, 16, 63, false)
	if err == nil {
		t.Errorf("BitString error")
	}

	v, bitlen, err = EncBitString(in, 100, 0, 63, false)
	if err == nil {
		t.Errorf("BitString error")
	}

	in = []uint8{0x00, 0x00, 0x00}
	v, bitlen, err = EncBitString(in, 25, 22, 32, false)
	if err == nil {
		t.Errorf("BitString error")
	}

	in = []uint8{0x00, 0x00}
	v, bitlen, err = EncBitString(in, 16, 64, 64, false)
	if err == nil {
		t.Errorf("BitString error")
	}

	//fixed length case. but not implemented yet.
	in = []uint8{0x00, 0x00}
	v, bitlen, err = EncBitString(in, 16, 16, 16, false)
	expect := []uint8{0x00, 0x00}
	if bitlen != 16 || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", 16, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	in = []uint8{0x00, 0x00}
	v, bitlen, err = EncBitString(in, 16, 0, 255, false)
	expect = []uint8{0x10, 0x00, 0x00}
	expectlen := 24
	if bitlen != expectlen || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", expectlen, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	in = []uint8{0x00, 0x00, 0x00}
	v, bitlen, err = EncBitString(in, 23, 22, 32, false)
	expect = []uint8{0x00, 0x80, 0x00, 0x00}
	expectlen = 27
	if bitlen != expectlen || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", expectlen, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	in = []uint8{0x00, 0x00, 0x00, 0x00}
	v, bitlen, err = EncBitString(in, 25, 22, 32, false)
	expect = []uint8{0x06, 0x00, 0x00, 0x00}
	expectlen = 29
	if bitlen != expectlen || compareSlice(expect, v) == false {
		t.Errorf("bitlen expect: %d, actual %d", expectlen, bitlen)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}
}

func TestOctetString(t *testing.T) {

	in := make([]uint8, 10, 10)
	pv, plen, v, err := EncOctetString(in, 16, 64, false)
	if err == nil {
		t.Errorf("BitString error")
	}

	min := 8
	max := 8
	in = make([]uint8, max, max)
	pexpect := []uint8{}
	expect := in
	pv, plen, v, err = EncOctetString(in, min, max, false)
	expectplen := 0
	if compareSlice(pexpect, pv) == false || plen != expectplen ||
		compareSlice(expect, v) == false {
		t.Errorf("plen expect: %d, actual %d", expectplen, plen)
		t.Errorf("value pexpect: 0x%02x, actual 0x%02x", pexpect, pv)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	min = 2
	max = 2
	in = make([]uint8, max, max)
	pexpect = []uint8{0x00, 0x00, 0x00}
	expect = []uint8{}
	pv, plen, v, err = EncOctetString(in, min, max, true)
	expectplen = 17
	if compareSlice(pexpect, pv) == false || plen != expectplen ||
		compareSlice(expect, v) == false {
		t.Errorf("plen expect: %d, actual %d", expectplen, plen)
		t.Errorf("value pexpect: 0x%02x, actual 0x%02x", pexpect, pv)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	min = 8
	max = 8
	in = make([]uint8, max, max)
	pexpect = []uint8{0x00}
	expect = in
	pv, plen, v, err = EncOctetString(in, min, max, true)
	expectplen = 1
	if compareSlice(pexpect, pv) == false || plen != expectplen ||
		compareSlice(expect, v) == false {
		t.Errorf("plen expect: %d, actual %d", expectplen, plen)
		t.Errorf("value pexpect: 0x%02x, actual 0x%02x", pexpect, pv)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}

	min = 0
	max = 7
	in = make([]uint8, 3, 3)
	pexpect = []uint8{0x03}
	expect = in
	pv, plen, v, err = EncOctetString(in, min, max, true)
	expectplen = 5
	if compareSlice(pexpect, pv) == false || plen != expectplen ||
		compareSlice(expect, v) == false {
		t.Errorf("plen expect: %d, actual %d", expectplen, plen)
		t.Errorf("value pexpect: 0x%02x, actual 0x%02x", pexpect, pv)
		t.Errorf("value expect: 0x%02x, actual 0x%02x", expect, v)
	}
}