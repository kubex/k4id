package k4id

import (
	"math/big"
	"strings"
	"time"
)

const checksumSize = 1
const defaultIDLength = 15
const defaultTimeGenerator = TimeGeneratorMicro

var globalIDHost Generator

type TimeGenerator int

const (
	//TimeGeneratorNano 11 Char string
	TimeGeneratorNano TimeGenerator = iota
	// TimeGeneratorMicro 9 Char string
	TimeGeneratorMicro
	// TimeGeneratorMilli 7-8 Char string
	TimeGeneratorMilli
	// TimeGeneratorSecond 6 Char string
	TimeGeneratorSecond
	// TimeGeneratorMinute 5 Char string
	TimeGeneratorMinute
	// TimeGeneratorHour 4 Char string
	TimeGeneratorHour
	// TimeGeneratorDay 3 Char string
	TimeGeneratorDay
)

// Generate returns a string representation of the time in the given format
func (t TimeGenerator) Generate(src time.Time) string {
	var i big.Int
	switch t {
	case TimeGeneratorNano:
		i.SetInt64(src.UnixNano())
	case TimeGeneratorMicro:
		i.SetInt64(src.UnixMicro())
	case TimeGeneratorMilli:
		i.SetInt64(src.UnixMilli())
	case TimeGeneratorSecond:
		i.SetInt64(src.Unix())
	case TimeGeneratorMinute:
		i.SetInt64(src.Unix() / 60)
	case TimeGeneratorHour:
		i.SetInt64(src.Unix() / 3600)
	case TimeGeneratorDay:
		i.SetInt64(src.Unix() / 86400)
	}
	return i.Text(62)
}

func (t TimeGenerator) Parse(src string) time.Time {
	var i big.Int
	i.SetString(src, 62)
	switch t {
	case TimeGeneratorNano:
		return time.Unix(0, i.Int64())
	case TimeGeneratorMicro:
		return time.Unix(0, i.Int64()*1000)
	case TimeGeneratorMilli:
		return time.Unix(0, i.Int64()*1000000)
	case TimeGeneratorSecond:
		return time.Unix(i.Int64(), 0)
	case TimeGeneratorMinute:
		return time.Unix(i.Int64()*60, 0)
	case TimeGeneratorHour:
		return time.Unix(i.Int64()*3600, 0)
	case TimeGeneratorDay:
		return time.Unix(i.Int64()*86400, 0)
	}
	return time.Time{}
}

func init() {
	globalIDHost = DefaultGenerator()
}

// Generator is a unique ID generator that can be configured
type Generator struct {
	hostID       string
	hostIDLength int
	idLength     int
	timeSize     TimeGenerator
	generation   chan bool
	withTime     *time.Time
}

// DefaultGenerator returns a new Generator with default configuration
func DefaultGenerator() Generator {
	h := Generator{
		idLength: defaultIDLength,
		timeSize: defaultTimeGenerator,
	}
	h.randomHostID()
	h.generation = make(chan bool, 1)
	return h
}

func NewGenerator(timeSize TimeGenerator) Generator {
	h := DefaultGenerator()
	h.SetTimeSize(timeSize)
	return h
}

// SetHostID sets the host ID to be used when generating IDs
func (h *Generator) SetHostID(id string) {
	h.hostID = id
	h.hostIDLength = len(h.hostID)
}

// GetHostID returns the current host ID
func (h *Generator) GetHostID() string { return h.hostID }
func (h *Generator) randomHostID()     { h.SetHostID(randomString(2)) }

// SetBaseLength sets the length of the base string
func (h *Generator) SetBaseLength(size int) { h.idLength = size }

// SetTimeSize sets the time size
func (h *Generator) SetTimeSize(size TimeGenerator) { h.timeSize = size }

// SetTime sets the time to be used when generating IDs, this may result in duplicate IDs being generated
func (h *Generator) SetTime(when time.Time) { h.withTime = &when }

// ClearTime clears the time to be used when generating IDs, using the current time for ID generation
func (h *Generator) ClearTime() { h.withTime = nil }

// New returns a new ID from the generator
func (h *Generator) New() ID {
	if h.hostID == "" {
		h.randomHostID()
	}

	h.generation <- true
	i := ID{}
	i.uniqueKey = h.randomID()
	i.verification = i.checkSum(i.uniqueKey)
	if h.timeSize == TimeGeneratorNano && h.idLength < 15 {
		//Sleep for a nanosecond to ensure uniqueness when precision is needed
		time.Sleep(time.Nanosecond)
	}
	<-h.generation
	return i
}

func (h *Generator) time() time.Time {
	if h.withTime == nil {
		return time.Now()
	}
	return *h.withTime
}

func (h *Generator) randomID() string {
	tId := h.reverse(h.timeSize.Generate(h.time()))
	useLen := h.idLength - len(tId) - len(h.hostID)
	return h.fixLen(tId+h.hostID+randomString(useLen), h.idLength)
}

func (h *Generator) reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (h *Generator) fixLen(input string, reqLen int) string {
	srcLen := len(input)
	if srcLen == reqLen {
		return input
	}
	if srcLen > reqLen {
		return input[:reqLen]
	}
	return input + strings.Repeat("X", reqLen-srcLen)
}

func (h *Generator) ExtractTime(fullID string) time.Time {
	id := FromString(fullID).uniqueKey
	timeSize := 3
	switch h.timeSize {
	case TimeGeneratorNano:
		timeSize = 11
	case TimeGeneratorMicro:
		timeSize = 9
	case TimeGeneratorMilli:
		timeSize = 7
	case TimeGeneratorSecond:
		timeSize = 6
	case TimeGeneratorMinute:
		timeSize = 5
	case TimeGeneratorHour:
		timeSize = 4
	}
	if len(id) < timeSize {
		return time.Time{}
	}

	return h.timeSize.Parse(h.reverse(id[:timeSize]))
}
