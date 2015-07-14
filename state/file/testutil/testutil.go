/*
Package testutil provides the utilty functions to test the file package.
*/
package testutil

import (
	"io/ioutil"
	"math/rand"
	"path"
)

// Directory represents a working directory in the test.
type Directory struct {
	dir string
}

// TempDir creates a temporary directory and returns the Directory struct.
func TempDir() *Directory {
	dir, err := ioutil.TempDir("", "file_test_")
	if err != nil {
		panic(err)
	}
	return &Directory{dir}
}

// MakeDummyFile creates the dummy file which contains a thousands of random
// ascii characters.
func (d *Directory) MakeDummyFile(prefix string) string {
	file, err := ioutil.TempFile(d.dir, prefix)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		n := rand.Intn(126-32) + 32
		file.Write([]byte{byte(n)})
	}
	return file.Name()
}

// NewFilePath returns the new file path for the named file.
func (d *Directory) NewFilePath(name string) string {
	return path.Join(d.dir, name)
}
