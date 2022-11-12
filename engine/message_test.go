package engine_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/umbe77/yasb/engine"
)

func TestParseMessage(t *testing.T) {
	input := `{
        "name": "Roberto",
        "age": 45,
        "birthDate": "1977-04-24T00:02:00Z",
        "float": 23.4,
        "b": false,
        "a": ["p", 10, true, {"a": "pluto"}, [true, {"b": 45}]],
        "sub": {
            "subName": "Ughi",
            "item":{
                "prova2": "pluto"
            }
        }
    }`

	inputMsg := make(map[string]any)
	if err := json.Unmarshal([]byte(input), &inputMsg); err != nil {
		t.Errorf("Cannot deserialize input message: %s", err.Error())
		return
	}
	out := engine.ParseInput(inputMsg)
	v, ok := out["name"]
	if !ok {
		t.Error("name property not in out Message")
	}
	val, ok := v.(string)
	if !ok {
		t.Errorf("name property values should be a string --> %q", v)
	}
	if val != "Roberto" {
		t.Errorf("Got %s, EXPECTED Roberto", val)
	}

	v, ok = out["age"]
	if !ok {
		t.Error("age property not in out Message")
	}
	valInt, ok := v.(int)
	if !ok {
		t.Errorf("age property values should be an int --> %q", v)
	}
	if valInt != 45 {
		t.Errorf("Got %d, EXPECTED 45", valInt)
	}

	v, ok = out["float"]
	if !ok {
		t.Error("float property not in out Message")
	}
	valfloat, ok := v.(float64)
	if !ok {
		t.Errorf("float property values should be a float64 --> %q", v)
	}
	if valfloat != 23.4 {
		t.Errorf("Got %.3f, EXPECTED 23.4", valfloat)
	}

	v, ok = out["b"]
	if !ok {
		t.Error("b property not in out Message")
	}
	valBool, ok := v.(bool)
	if !ok {
		t.Errorf("b property values should be a bool --> %q", v)
	}
	if valBool { //valBool should be false
		t.Errorf("Got %t, EXPECTED false", valBool)
	}

	v, ok = out["sub->subName"]
	if !ok {
		t.Error("sub->subName property not in out Message")
	}
	valSubName, ok := v.(string)
	if !ok {
		t.Errorf("sub->subName property values should be a string --> %q", v)
	}
	if valSubName != "Ughi" {
		t.Errorf("Got %s, EXPECTED Ughi", valSubName)
	}

	v, ok = out["sub->item->prova2"]
	if !ok {
		t.Error("sub->item->prova2 property not in out Message")
	}
	valProva2, ok := v.(string)
	if !ok {
		t.Errorf("sub->item->prova2 property values should be a string --> %q", v)
	}
	if valProva2 != "pluto" {
		t.Errorf("Got %s, EXPECTED pluto", valProva2)
	}

	v, ok = out["birthDate"]
	if !ok {
		t.Error("birthDate property not in out Message")
	}
	valTime, ok := v.(time.Time)
	if !ok {
		t.Errorf("birthDate property values should be a time --> %q", v)
	}
	expectedTime := time.Date(1977, time.April, 24, 0, 2, 0, 0, time.UTC)
	if !valTime.Equal(expectedTime) {
		t.Errorf("Got %s, EXPECTED 24/04/1977 00:02:00", valProva2)
	}

	v, ok = out["a"]
	if !ok {
		t.Error("a property not in out Message")
	}
	valSlice, ok := v.([]any)
	if !ok {
		t.Errorf("a property values should be an array --> %q", v)
	}

	if len(valSlice) != 5 {
		t.Errorf("a property values should have a length of 5 elements got  --> %d", len(valSlice))
	}
	for i, item := range valSlice {
		switch i {
		case 0:
			val, ok := item.(string)
			if !ok {
				t.Errorf("a[0] property values should be a string --> %q", v)

			}
			if val != "p" {
				t.Errorf("Got %s, EXPECTED p", val)
			}
		case 1:
			val, ok := item.(int)
			if !ok {
				t.Errorf("a[1] property values should be a int --> %q", v)

			}
			if val != 10 {
				t.Errorf("Got %d, EXPECTED p", val)
			}
		case 2:
			val, ok := item.(bool)
			if !ok {
				t.Errorf("a[2] property values should be a bool --> %q", v)
			}
			if !val {
				t.Errorf("Got %t, EXPECTED p", val)
			}
		case 3:
			v, ok := item.(engine.TaskMessage)
			if !ok {
				t.Errorf("a[3] property values should be a TaskMessage --> %q", v)
			}
			subItem, ok := v["a"]
			if !ok {
				t.Error("a[3][a] property should be in message")
			}
			val, ok := subItem.(string)
			if !ok {
				t.Errorf("a[3][a] property values should be a string --> %q", v)
			}
			if val != "pluto" {
				t.Errorf("Got %s, EXPECTED p", val)
			}
		case 4:
			val, ok := item.([]any)
			if !ok {
				t.Errorf("a[4] should be a Slice --> %q", val)
			}
			if len(val) != 2 {
				t.Errorf("a[4] should have 2 items, got: %d", len(val))
			}
			for index, subArrayItem := range val {
				switch index {
				case 0:
					val, ok := subArrayItem.(bool)
					if !ok {
						t.Errorf("a[4][0] property values should be a bool --> %q", v)
					}
					if !val {
						t.Errorf("Got %t, EXPECTED true", val)
					}
				case 1:
					v, ok := subArrayItem.(engine.TaskMessage)
					if !ok {
						t.Errorf("a[4][1] property values should be a TaskMessage --> %q", v)
					}
					subItem, ok := v["b"]
					if !ok {
						t.Error("a[4][1][b] property should be in message")
					}
					val, ok := subItem.(int)
					if !ok {
						t.Errorf("a[4][1[b] property values should be a int --> %q", v)
					}
					if val != 45 {
						t.Errorf("Got %d, EXPECTED 45", val)
					}
				}
			}
		}

	}
}
