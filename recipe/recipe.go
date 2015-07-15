package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/resource"
)

// Recipe represents the recipe which desribes disired states of resources.
type Recipe struct {
	Config    map[string]string
	Resources []resource.Resource
}

var fileNames = []string{
	"recipe.json",
	"recipe.rb",
	"recipe.yaml",
	"recipe.yml",
}

// ReadRecipe reads the named recipe file. If the name is not specified, it
// searchs the recipe file in specified directory.
// If failed to read recipe, it returns an error.
func ReadRecipe(name string, dir string) (*Recipe, error) {
	var r io.Reader
	var err error
	if name != "" {
		name, r, err = openFile(name)
	} else {
		name, r, err = findFile(dir)
	}
	if err != nil {
		return nil, err
	}

	logger.Debugf("------ reading recipe file: %s", name)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Fatalf("can not read the file: %s", err)
	}

	switch {
	case strings.HasSuffix(name, ".json"):
		r, err := ParseJSON(data)
		if err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				line, pos := getCaretPos(data, int(e.Offset))
				return nil, fmt.Errorf("can not parse the JSON file: %s near line %d, pos %d", err, line, pos)
			}
			return nil, fmt.Errorf("can not parse the JSON file: %s", err)
		}
		return r, nil
	}

	return nil, fmt.Errorf("unsupported file format")
}

func openFile(name string) (path string, r io.Reader, err error) {
	path, err = filepath.Abs(name)
	if err != nil {
		return path, nil, err
	}
	r, err = os.Open(name)
	return
}

func findFile(dir string) (name string, r io.Reader, err error) {
	var file *os.File
	for _, name = range fileNames {
		file, err = os.Open(path.Join(dir, name))
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return name, file, err
		}
	}
	if file == nil {
		err = errors.New("the recipe file is not found")
	}
	return name, file, err
}

func getCaretPos(data []byte, off int) (line int, pos int) {
	line = 1
	pos = 0
	for i, b := range data {
		if i >= off {
			break
		}
		if b == '\n' {
			line++
			pos = 0
			continue
		}
		pos++
	}
	return
}
