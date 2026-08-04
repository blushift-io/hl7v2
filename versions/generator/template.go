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

type Template struct {
	typ    SpecType
	t      *template.Template
	output string
}

func NewTemplate(typ SpecType, t *template.Template, output string) *Template {
	return &Template{
		typ:    typ,
		t:      t,
		output: output,
	}
}

func ParseTemplate(typ SpecType, name, content, output string) (*Template, error) {
	t, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, err
	}

	return NewTemplate(typ, t, name), nil
}

func ReadTemplate(typ SpecType, path, output string) (*Template, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return NewTemplate(typ, t, name), nil
}

func ReadTemplateFS(typ SpecType, fs fs.FS, output string) (*Template, error) {
	t, err := template.ParseFS(fs, output)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(output)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return NewTemplate(typ, t, name), nil
}
