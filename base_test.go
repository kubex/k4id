package k4id

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"testing"
)

// Write a test for encodeB63
func TestEncodeB63(t *testing.T) {
	decoded := "aRS3gQXoI7UvTxYX"
	x, e := Base63.Decode(decoded)
	log.Println(x)
	log.Println(fmt.Sprintf("%x", x), e)
	encoded := Base63.Encode(x)
	log.Println(encoded)
	encoded2 := Base63.Encode([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})
	log.Println(encoded2)
}

func TestBase36(t *testing.T) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, 1234567890)
	log.Println(Base36.Encode(bs))
	log.Println(Base36.Decode("KF12OI"))

	bs2 := make([]byte, 4)
	binary.BigEndian.PutUint32(bs2, 2147483646)
	log.Println(Base36.Encode(bs2))
	log.Println(Base36.Decode("ZIK0ZI"))

	bs64 := make([]byte, 8)
	binary.BigEndian.PutUint64(bs64, 9223372036854775807)
	log.Println(Base36.Encode(bs64))
	log.Println(Base36.Decode("1Y2P0IJ32E8E7"))
}

func TestEncodeUInt64(t *testing.T) {
	maxB36 := Base36.EncodeUInt64(math.MaxInt64)
	if maxB36 != "1Y2P0IJ32E8E7" {
		t.Errorf("Expected 1Y2P0IJ32E8E7 but got %s", maxB36)
	}

	maxiB36 := Base36.EncodeUInt64(math.MaxUint64)
	if maxiB36 != "3W5E11264SGSF" {
		t.Errorf("Expected 3W5E11264SGSF but got %s", maxiB36)
	}
}

func TestDecodeUInt64(t *testing.T) {
	maxB36 := Base36.DecodeToUInt64("1Y2P0IJ32E8E7")
	if maxB36 != math.MaxInt64 {
		t.Errorf("Expected MaxInt64 but got %d", maxB36)
	}
	maxiB36 := Base36.DecodeToUInt64("3W5E11264SGSF")
	if maxiB36 != math.MaxUint64 {
		t.Errorf("Expected MaxInt64 but got %d", maxiB36)
	}
}

func TestUInt64(t *testing.T) {
	for i := uint64(0); i < math.MaxInt16; i++ {
		encoded := Base36.EncodeUInt64(i)
		decoded := Base36.DecodeToUInt64(encoded)
		if decoded != i {
			t.Errorf("Expected %d but got %d for encoded %s", i, decoded, encoded)
		}
	}
}

func TestCompact62(t *testing.T) {
	base := Base36
	log.Println("Result", base.CompactHex(""))
	log.Println("Result", base.CompactHex("0"))
	log.Println("Result", base.CompactHex("1"))
	log.Println("Result", base.CompactHex("2"))
	log.Println("Result", base.CompactHex("3"))
	log.Println("Result", base.CompactHex("4"))
	log.Println("Result", base.CompactHex("5"))
	log.Println("Result", base.CompactHex("6"))
	log.Println("Result", base.CompactHex("7"))
	log.Println("Result", base.CompactHex("8"))
	log.Println("Result", base.CompactHex("9"))
	log.Println("Result", base.CompactHex("a"))
	log.Println("Result", base.CompactHex("A"))
	log.Println("Result", base.CompactHex("b"))
	log.Println("Result", base.CompactHex("B"))
	log.Println("Result", base.CompactHex("c"))
	log.Println("Result", base.CompactHex("d"))
	log.Println("Result", base.CompactHex("e"))
	log.Println("Result", base.CompactHex("f"))
	log.Println("Result", base.CompactHex("10"))
	log.Println("Result", base.CompactHex("13"))
	log.Println("Result", base.CompactHex("ffff"))
	log.Println("Result", base.CompactHex("ffffffff"))
	log.Println("Result", base.CompactHex("ffffffffffff"))
	log.Println("Result", base.CompactHex("ffffffffffffffff"))
	log.Println("Result", base.CompactHex("ffffffffffffffffffff"))
	log.Println("Result", base.CompactHex("ffffffffffffffffffffffff"))                                         //LygHa16AHYF4gfFC3
	log.Println("Result", base.CompactHex("ffffffffffffffffffffffffffff"))                                     //LygHa16AHYF1HvWXNAa7
	log.Println("Result", base.CompactHex("fffffffffffffffffffffffffffff"))                                    //LygHa16AHYF1HvWXNAa7
	log.Println("Result", base.CompactHex("ffffffffffffffffffffffffffffffff"))                                 //LygHa16AHYFLygHa16AHYF
	log.Println("Result", base.CompactHex("0000000000000000000000000000000000000000000000000000000000000000")) //LygHa16AHYFLygHa16AHYF
	log.Println("Result", base.CompactHex("100000000000000000000000010100000000000000000000"))                 //LygHa16AHYFLygHa16AHYF
	log.Println("Result", base.CompactHex("424663d61809925fa407c54ae4278d7f"))                                 //LygHa16AHYFLygHa16AHYF
}
