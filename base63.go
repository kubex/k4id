package k4id

import (
	"bytes"
	"errors"
)

const baseNum = 63

var base63Alpha = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz?")
var base63Map = make(map[rune]int)

func init() {
	for i := 0; i < baseNum; i++ {
		base63Map[base63Alpha[i]] = i
	}
}

func encodeB63(source []byte) string {
	if len(source) == 0 {
		return ""
	}

	digits := []int{0}

	for i := 0; i < len(source); i++ {
		carry := int(source[i])

		for j := 0; j < len(digits); j++ {
			carry += digits[j] << 8
			digits[j] = carry % baseNum
			carry = carry / baseNum
		}

		for carry > 0 {
			digits = append(digits, carry%baseNum)
			carry = carry / baseNum
		}
	}

	var res bytes.Buffer

	for k := 0; source[k] == 0 && k < len(source)-1; k++ {
		res.WriteRune(base63Alpha[0])
	}

	for q := len(digits) - 1; q >= 0; q-- {
		res.WriteRune(base63Alpha[digits[q]])
	}

	return res.String()
}

func decodeB63(source string) ([]byte, error) {
	if len(source) == 0 {
		return []byte{}, nil
	}

	runes := []rune(source)

	byts := []byte{0}
	for i := 0; i < len(runes); i++ {
		value, ok := base63Map[runes[i]]

		if !ok {
			return nil, errors.New("non base62 character")
		}

		carry := value

		for j := 0; j < len(byts); j++ {
			carry += int(byts[j]) * baseNum
			byts[j] = byte(carry & 0xff)
			carry >>= 8
		}

		for carry > 0 {
			byts = append(byts, byte(carry&0xff))
			carry >>= 8
		}
	}

	for k := 0; runes[k] == base63Alpha[0] && k < len(runes)-1; k++ {
		byts = append(byts, 0)
	}

	// Reverse bytes
	for i, j := 0, len(byts)-1; i < j; i, j = i+1, j-1 {
		byts[i], byts[j] = byts[j], byts[i]
	}

	return byts, nil
}
