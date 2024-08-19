package k4id

import (
	"encoding/binary"
	"fmt"
	"log"
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
	binary.LittleEndian.PutUint32(bs, 1234567890)
	log.Println(Base36.Encode(bs))
	log.Println(Base36.Decode("KF12OI"))

	bs2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs2, 2147483646)
	log.Println(Base36.Encode(bs2))
	log.Println(Base36.Decode("ZIK0ZI"))

	bs64 := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs64, 992147483646)
	log.Println(Base36.Encode(bs64))
	log.Println(Base36.Decode("CNSAZTVI"))
}
