package engine

import (
	"fmt"
	"time"

	"github.com/umbe77/yasb/stack"
)

type TaskMessage map[string]any

type Message struct {
	Outputs map[string]TaskMessage
}

const rootPath string = ""

func ParseInput(in map[string]any) TaskMessage {
	return parseInputInternal(in, rootPath)
}

func getPrimaryValue(v any) any {
	switch v := v.(type) {
	case float64:
		x := int(v)
		if v-float64(x) == 0 {
			return x
		}
		return v
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return t
		}
		return v
	default:
		return v
	}
}

func getArrayValue(in []any) []any {
	out := make([]any, 0, len(in))
	for _, v := range in {
		switch v := v.(type) {
		case []any:
			out = append(out, getArrayValue(v))
		case map[string]any:
			out = append(out, parseInputInternal(v, rootPath))
		default:
			out = append(out, getPrimaryValue(v))
		}
	}
	return out
}

type messageVisit struct {
	m    TaskMessage
	path string
}

func formatPath(path, key string) string {
	if path != "" {
		return fmt.Sprintf("%s->%s", path, key)
	}
	return key
}

func parseInputInternal(in map[string]any, path string) TaskMessage {
	out := make(TaskMessage)

	s := stack.New[messageVisit]()
	s.Push(messageVisit{
		in,
		path,
	})

	for !s.IsEmpty() {
		msg := s.Pop()
		for k, v := range msg.m {
			switch v := v.(type) {
			case map[string]any:
				s.Push(messageVisit{
					v,
					formatPath(msg.path, k),
				})
			case []any:
				out[formatPath(msg.path, k)] = getArrayValue(v)
			default:
				out[formatPath(msg.path, k)] = getPrimaryValue(v)
			}
		}
	}

	return out
}
