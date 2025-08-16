package json_test

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	data := make(map[string]interface{})
	data["name"] = "test"
	data["age"] = 18
	data["address"] = "test"
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	if bb, err := json.Marshal(string(b)); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(bb))
	}
}

func TestJson2(t *testing.T) {
	a := `"{\"address\":\"test\",\"age\":18,\"name\":\"test\"}"`
	var raw string
	if err := json.Unmarshal([]byte(a), &raw); err != nil {
		t.Fatal(err)
	}
	t.Log(raw)
}
