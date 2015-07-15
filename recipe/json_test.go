package recipe

import "testing"

var input = []byte(`{
  "resources": [
    {
      "type": "file",
      "path": "/tmp/file"
    }
  ]
}`)

func TestParseJSON(t *testing.T) {
	recipe, err := ParseJSON(input)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if len(recipe.Resources) != 1 {
		t.Errorf("parsed 1 resource expected, but got %d", len(recipe.Resources))
	}
}
