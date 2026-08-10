package generator

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates
var defaultTemplates embed.FS

// Template wraps a parsed text template associated with a SpecType.
type Template struct {
	typ    SpecType
	t      *template.Template
	output string
}

// NewTemplate creates a new Template instance for a SpecType and output name.
func NewTemplate(typ SpecType, t *template.Template, output string) *Template {
	return &Template{
		typ:    typ,
		t:      t,
		output: output,
	}
}

// ParseTemplate parses a template content string and returns a new Template.
func ParseTemplate(typ SpecType, name, content, output string) (*Template, error) {
	t, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, err
	}

	return NewTemplate(typ, t, name), nil
}

// ReadTemplate reads and parses a template file from disk.
func ReadTemplate(typ SpecType, path, output string) (*Template, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return NewTemplate(typ, t, name), nil
}

// ReadTemplateFS reads and parses a template file from a filesystem.
func ReadTemplateFS(typ SpecType, fs fs.FS, output string) (*Template, error) {
	t, err := template.ParseFS(fs, output)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(output)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return NewTemplate(typ, t, name), nil
}
