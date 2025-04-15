package k4id

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

func Hashed(input string, length uint8) string {
	h := make([]byte, length+10)
	sha3.ShakeSum256(h, []byte(input))
	return Base62.CompactHex(hex.EncodeToString(h))[:length]
}
