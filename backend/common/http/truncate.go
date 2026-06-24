package http

import (
	"encoding/json"
	"fmt"
)

const (
	maxBodyLogBytes  = 512
	maxArrayLogItems = 6
)

func truncateBodyForLog(body string) string {
	if len(body) == 0 {
		return ""
	}
	if len(body) <= maxBodyLogBytes {
		return body
	}

	var parsed any
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return headTailTruncate(body)
	}

	if arr, ok := parsed.([]any); ok && len(arr) > maxArrayLogItems {
		parsed = truncateArray(arr)
	}
	truncateJSONArrays(parsed)

	out, err := json.Marshal(parsed)
	if err != nil {
		return headTailTruncate(body)
	}
	return string(out)
}

func truncateJSONArrays(v any) {
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			if arr, ok := child.([]any); ok && len(arr) > maxArrayLogItems {
				val[k] = truncateArray(arr)
				truncateJSONArrays(val[k])
			} else {
				truncateJSONArrays(child)
			}
		}
	case []any:
		for _, item := range val {
			truncateJSONArrays(item)
		}
	}
}

func truncateArray(arr []any) []any {
	keep := maxArrayLogItems - 1
	truncated := len(arr) - keep - 1
	result := make([]any, 0, keep+2)
	result = append(result, arr[:keep]...)
	result = append(result, map[string]any{"...": fmt.Sprintf("%d items truncated", truncated)})
	result = append(result, arr[len(arr)-1])
	return result
}

func headTailTruncate(s string) string {
	half := maxBodyLogBytes / 2
	truncated := len(s) - 2*half
	return s[:half] + fmt.Sprintf(" ... (%d bytes truncated) ... ", truncated) + s[len(s)-half:]
}
