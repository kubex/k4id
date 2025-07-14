package k4id

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

var Base16 = NewBase([]rune("0123456789ABCDEF"))
var Base32 = NewBase([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"))
var Base36 = NewBase([]rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
var Base62 = NewBase([]rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"))
var Base63 = NewBase([]rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz?"))

func NewBase(alpha []rune) *Base {
	b := &Base{
		base:    len(alpha),
		alpha:   alpha,
		baseMap: make(map[rune]int),
	}

	for i := 0; i < b.base; i++ {
		b.baseMap[b.alpha[i]] = i
	}
	return b
}

type Base struct {
	base    int
	alpha   []rune
	baseMap map[rune]int
}

func (b Base) Encode(source []byte) string {
	source = bytes.Trim(source, "\x00")
	return b.EncodeWithTrailing(source)
}

func (b Base) EncodeWithTrailing(source []byte) string {
	if len(source) == 0 {
		return ""
	}

	digits := []int{0}

	for i := 0; i < len(source); i++ {
		carry := int(source[i])

		for j := 0; j < len(digits); j++ {
			carry += digits[j] << 8
			digits[j] = carry % b.base
			carry = carry / b.base
		}

		for carry > 0 {
			digits = append(digits, carry%b.base)
			carry = carry / b.base
		}
	}

	var res bytes.Buffer

	for k := 0; source[k] == 0 && k < len(source)-1; k++ {
		res.WriteRune(b.alpha[0])
	}

	for q := len(digits) - 1; q >= 0; q-- {
		res.WriteRune(b.alpha[digits[q]])
	}

	return res.String()
}

func (b Base) Decode(source string) ([]byte, error) {
	if len(source) == 0 {
		return []byte{}, nil
	}

	runes := []rune(source)

	byts := []byte{0}
	for i := 0; i < len(runes); i++ {
		value, ok := b.baseMap[runes[i]]

		if !ok {
			return nil, errors.New("non base62 character")
		}

		carry := value

		for j := 0; j < len(byts); j++ {
			carry += int(byts[j]) * b.base
			byts[j] = byte(carry & 0xff)
			carry >>= 8
		}

		for carry > 0 {
			byts = append(byts, byte(carry&0xff))
			carry >>= 8
		}
	}

	for k := 0; runes[k] == b.alpha[0] && k < len(runes)-1; k++ {
		byts = append(byts, 0)
	}

	// Reverse bytes
	for i, j := 0, len(byts)-1; i < j; i, j = i+1, j-1 {
		byts[i], byts[j] = byts[j], byts[i]
	}

	return byts, nil
}

func (b Base) EncodeUInt64(source uint64) string {
	bs64 := make([]byte, 8)
	binary.BigEndian.PutUint64(bs64, source)
	return b.EncodeWithTrailing(bs64)
}

func (b Base) DecodeToUInt64(source string) uint64 {
	decoded, err := b.Decode(source)
	if err != nil {
		return 0
	}
	decodedB := make([]byte, 8)
	copy(decodedB[8-len(decoded):], decoded)
	return binary.BigEndian.Uint64(decodedB)
}

func (b Base) CompactHex(source string) string {
	final := ""
	chunks := hexChunks(source)
	for _, chunk := range chunks {
		if chunk == 0 {
			final += "0"
		} else {
			final += b.EncodeUInt64(chunk)
		}
	}
	return final
}

func hexChunks(source string) []uint64 {
	chunks := make([]uint64, 0, len(source)/16)
	for i := 0; i < len(source); i += 16 {
		end := i + 16
		if end > len(source) {
			end = len(source)
		}
		chunk, err := strconv.ParseUint(source[i:end], 16, 64)
		if err != nil {
			return nil
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}
