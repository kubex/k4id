package k4id

import (
	"strconv"
	"testing"
	"time"
)

func Test_Hashed(t *testing.T) {
	conflictMap := map[string]string{}
	for i := 0; i < 10000; i++ {
		hashVal := strconv.FormatInt(time.Now().UnixNano(), 10)
		hash := Hashed(hashVal, 19).String()
		if val, ok := conflictMap[hash]; ok {
			t.Errorf("conflictMap exists: %v - %s - %s", hash, hashVal, val)
			return
		}
		conflictMap[hash] = hashVal
		time.Sleep(time.Nanosecond)
	}
}
