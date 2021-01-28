package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// OddPairsKey is key marker for not even pair
const oddPairsKey = "ODD_KEYVALS"

// pair ...
type pair struct {
	Key   string
	Value interface{}
}

// pairsList ...
type pairsList []pair

// NewPairsList ...
func newPairsList(pairs ...interface{}) pairsList {
	if len(pairs) == 0 {
		return nil
	}

	if len(pairs)%2 != 0 { // prevent odd pairs in log
		pairs = []interface{}{oddPairsKey, fmt.Sprint(pairs)}
	}

	n := len(pairs)
	result := make(pairsList, 0, n/2)
	for i := 1; i < n; i += 2 {
		result = append(result, pair{
			Key:   toString(pairs[i-1]),
			Value: pairs[i],
		})
	}

	return result
}

// String ...
func (pl pairsList) String() string {
	buf := bytes.NewBufferString("{")
	first := true
	for _, entry := range pl {
		if !first {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("%s: %v", entry.Key, entry.Value))
		first = false
	}
	_, _ = buf.WriteString("}")

	return buf.String()
}

// JSON ...
func (pl pairsList) JSON() string {
	result := make(map[string]interface{})
	for i, entry := range pl {
		k := entry.Key
		if _, ok := result[k]; ok {
			k = fmt.Sprintf("%s_%d", k, i)
		}
		result[k] = entry.Value
	}

	buf, err := json.Marshal(result)
	if err != nil {
		return "{}"
	}
	return string(buf)
}

// Flattened ...
func (pl pairsList) Flattened() []interface{} {
	result := make([]interface{}, 0, len(pl)*2)
	for _, entry := range pl {
		result = append(result, entry.Key, entry.Value)
	}
	return result
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case fmt.Stringer:
		return typed.String()
	case string:
		return typed
	default:
		return fmt.Sprint(value)
	}
}
