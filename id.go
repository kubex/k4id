//package:k4id is a package that provides a unique identifier generator.
package k4id

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func randomString(n int) string {
	if n < 1 {
		return ""
	}
	b := make([]byte, n+5)
	_, _ = rand.Read(b)
	return strings.ReplaceAll(strings.ReplaceAll(base64.RawStdEncoding.EncodeToString(b), "+", ""), "/", "")[:n]
}

// New creates a new random ID from the globalIDHost
func New() ID { return globalIDHost.New() }

// FromString hydrates an ID from a string
func FromString(input string) ID {
	i := ID{}
	if len(input) < checksumSize {
		return i
	}
	i.verification = input[:checksumSize]
	i.uniqueKey = input[checksumSize:]
	return i
}

// FromUUID converts a UUID to an ID
func FromUUID(input string) ID {
	rawID, _ := hex.DecodeString(strings.ReplaceAll(input, "-", ""))
	return FromString(strings.ReplaceAll(encodeB63(rawID), "?", ""))
}

// ID is a unique identifier, along with a verification checksum
type ID struct {
	uniqueKey    string
	verification string
}

// String returns the ID as a string
func (i ID) String() string {
	return strings.TrimSpace(i.verification + i.uniqueKey)
}

// UUID returns the ID as a UUID
func (i ID) UUID() string {
	if len(i.String()) > 21 {
		// IDs generated that are over 21 characters long, have the potential to be invalid uuids
		return ""
	}

	var rawID []byte
	var uuid string
	for xi := 20; xi < 24; xi++ {
		if len(uuid) < 32 {
			rawID, _ = decodeB63((i.String() + "??????????????????????")[:xi])
			uuid = hex.EncodeToString(rawID)
		}
	}

	return fmt.Sprintf("%s-%s-%s-%s-%s", uuid[:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])
}

// IsValid returns true if the ID has a valid verification checksum
func (i ID) IsValid() bool {
	if i.verification == "" || len(i.verification) != checksumSize {
		return false
	}

	return i.verification == i.checkSum(i.uniqueKey)
}

func (i ID) checkSum(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))[:checksumSize]
}
