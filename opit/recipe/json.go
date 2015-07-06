package recipe

import (
	"encoding/json"
	"fmt"

	"github.com/harukasan/orchestra-pit/opit/resource"
)

// ParseJSON parses the recipe file serialized by JSON.
func ParseJSON(data []byte) (*Recipe, error) {
	// unmarshal only the root node of the JSON.
	var root struct {
		Config    map[string]string
		Resources []json.RawMessage `json:"resources"`
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	recipe := &Recipe{
		Config: root.Config,
	}
	for _, r := range root.Resources {
		// unmarshal only the type attribute
		var attr struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(r, &attr); err != nil {
			return nil, err
		}

		res, err := unmarshalResource(r, attr.Type)
		if err != nil {
			return nil, err
		}
		recipe.Resources = append(recipe.Resources, res)
	}
	return recipe, nil
}

func unmarshalResource(j json.RawMessage, t string) (resource.Resource, error) {
	res := resource.New(t)
	if res == nil {
		return nil, fmt.Errorf("unknwon resource type: %s", t)
	}
	if err := json.Unmarshal(j, res); err != nil {
		return nil, err
	}
	return res, nil
}
