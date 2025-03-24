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
}
