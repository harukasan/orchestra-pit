package resource

import "github.com/harukasan/orchestra-pit/opit/resource/file"

type Resource interface {
	Apply() error
	Test() error
}

func New(t string) Resource {
	switch t {
	case "file":
		return &file.Resource{}
	}
	return nil
}
