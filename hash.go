package k4id

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

func Hashed(input string, length uint8) ID {
	h := make([]byte, length+10)
	sha3.ShakeSum256(h, []byte(input))
	i := ID{uniqueKey: Base62.CompactHex(hex.EncodeToString(h))[:length-checksumSize]}
	i.verification = i.checkSum(i.uniqueKey)
	return i
}
