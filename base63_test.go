package k4id

import (
	"fmt"
	"log"
	"testing"
)

// Write a test for encodeB63
func TestEncodeB62(t *testing.T) {
	decoded := "aRS3gQXoI7UvTxYX"
	x, e := decodeB63(decoded)
	log.Println(x)
	log.Println(fmt.Sprintf("%x", x), e)
	encoded := encodeB63(x)
	log.Println(encoded)
	encoded2 := encodeB63([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})
	log.Println(encoded2)
}
