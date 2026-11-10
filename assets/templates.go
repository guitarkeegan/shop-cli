package templates

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed templates
var raw embed.FS

var Templates = func() fs.FS {
	sub, err := fs.Sub(raw, "templates")
	if err != nil {
		log.Fatalf("could not sub templates: %v", err)
	}
	return sub
}()
