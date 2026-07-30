package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/schema"
	_ "github.com/blushift-io/hl7v2/schema/spec/v21"
	_ "github.com/blushift-io/hl7v2/schema/spec/v22"
	_ "github.com/blushift-io/hl7v2/schema/spec/v23"
	_ "github.com/blushift-io/hl7v2/schema/spec/v231"
	_ "github.com/blushift-io/hl7v2/schema/spec/v24"
	_ "github.com/blushift-io/hl7v2/schema/spec/v25"
	_ "github.com/blushift-io/hl7v2/schema/spec/v251"
	_ "github.com/blushift-io/hl7v2/schema/spec/v26"
	_ "github.com/blushift-io/hl7v2/schema/spec/v27"
	_ "github.com/blushift-io/hl7v2/schema/spec/v271"
	_ "github.com/blushift-io/hl7v2/schema/spec/v28"
)

//go:embed templates/*.gotmpl
var embeddedTemplates embed.FS

type SpecType int

const (
	SpecTypeUnknown SpecType = iota
	SpecTypeDatatype
	SpecTypeTable
	SpecTypeSegment
	SpecTypeMessage
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

var goKeywords = map[string]bool{
	"break": true, "default": true, "func": true, "interface": true, "select": true,
	"case": true, "defer": true, "go": true, "map": true, "struct": true,
	"chan": true, "else": true, "goto": true, "package": true, "switch": true,
	"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
	"continue": true, "for": true, "import": true, "return": true, "var": true,
}

type Deduper struct {
	used map[string]int
}

func NewDeduper() *Deduper {
	return &Deduper{used: make(map[string]int)}
}

func (d *Deduper) Name(raw string) string {
	base := SanitizeIdentifier(raw)
	d.used[base]++
	count := d.used[base]
	if count == 1 {
		return base
	}
	return fmt.Sprintf("%s%d", base, count)
}

func SanitizeIdentifier(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "Field"
	}

	words := nonAlphaNum.Split(s, -1)
	var sb strings.Builder
	for _, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		sb.WriteString(string(runes))
	}

	res := sb.String()
	if res == "" {
		return "Field"
	}

	if unicode.IsDigit([]rune(res)[0]) {
		res = "X" + res
	}

	if goKeywords[strings.ToLower(res)] {
		res = res + "_"
	}

	return res
}

func CommentLines(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = "// " + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func PackageName(ver string) string {
	ver = strings.TrimPrefix(ver, "v")
	ver = strings.ReplaceAll(ver, ".", "")
	return "v" + ver
}

type Generator struct {
	ver    hl7v2.Version
	schema *schema.Schema
	opts   *Options
}

func New(ver string, opts ...Option) (*Generator, error) {
	v := hl7v2.GetVersion(ver)
	if v == hl7v2.VersionUnknown {
		return nil, fmt.Errorf("unknown version %s", ver)
	}

	sch := schema.Open(v.String())
	if sch == nil {
		return nil, fmt.Errorf("unable to open schema for version %s", v.String())
	}

	options := NewOptions(opts...)

	return &Generator{
		ver:    v,
		schema: sch,
		opts:   options,
	}, nil
}

type DatatypeFieldModel struct {
	GoName string
	GoType string
	Tag    string
}

type DatatypeModel struct {
	ID          string
	Name        string
	Description string
	IsPrimitive bool
	Fields      []DatatypeFieldModel
}

type DatatypesFileData struct {
	PackageName string
	Datatypes   []DatatypeModel
}

type SegmentFieldModel struct {
	GoName string
	GoType string
	Tag    string
}

type SegmentModel struct {
	ID          string
	Name        string
	Description string
	Fields      []SegmentFieldModel
}

type SegmentsFileData struct {
	PackageName string
	Segments    []SegmentModel
}

type MessageFieldModel struct {
	GoName string
	GoType string
	Tag    string
}

type GroupModel struct {
	StructName string
	Fields     []MessageFieldModel
}

type MessageModel struct {
	ID          string
	Name        string
	Description string
	Groups      []GroupModel
	Fields      []MessageFieldModel
}

type MessagesFileData struct {
	PackageName string
	Messages    []MessageModel
}

type RegistryFileData struct {
	PackageName string
	VersionEnum string
}

type MessageTestsFileData struct {
	PackageName string
	Version     string
}

func (g *Generator) GenerateAll(outputBase string) error {
	pkgName := PackageName(g.ver.String())
	targetDir := filepath.Join(outputBase, pkgName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	if err := g.GenerateDatatypes(targetDir, pkgName); err != nil {
		return fmt.Errorf("datatypes gen error: %w", err)
	}
	if err := g.GenerateSegments(targetDir, pkgName); err != nil {
		return fmt.Errorf("segments gen error: %w", err)
	}
	if err := g.GenerateMessages(targetDir, pkgName); err != nil {
		return fmt.Errorf("messages gen error: %w", err)
	}
	if err := g.GenerateRegistry(targetDir, pkgName); err != nil {
		return fmt.Errorf("registry gen error: %w", err)
	}
	if err := g.GenerateMessageTests(targetDir, pkgName); err != nil {
		return fmt.Errorf("message tests gen error: %w", err)
	}

	return nil
}

func (g *Generator) GenerateDatatypes(targetDir, pkgName string) error {
	var models []DatatypeModel
	for _, dt := range g.schema.DataTypes() {
		m := DatatypeModel{
			ID:          dt.ID,
			Name:        dt.Name,
			Description: CommentLines(dt.Description),
			IsPrimitive: dt.IsPrimitive(),
		}

		if !dt.IsPrimitive() {
			d := NewDeduper()
			for idx, f := range dt.Fields {
				pos := idx + 1
				goName := d.Name(f.Name)
				if f.Name == "" {
					goName = d.Name(f.ID)
				}
				fType := f.DataType
				if fType == "" {
					fType = "ST"
				}
				if f.Repeatable() {
					fType = "[]" + fType
				}
				m.Fields = append(m.Fields, DatatypeFieldModel{
					GoName: goName,
					GoType: fType,
					Tag:    fmt.Sprintf("%d", pos),
				})
			}
		}
		models = append(models, m)
	}

	data := DatatypesFileData{
		PackageName: pkgName,
		Datatypes:   models,
	}

	return renderAndWrite("templates/datatypes.gotmpl", data, filepath.Join(targetDir, "datatypes.go"))
}

func (g *Generator) GenerateSegments(targetDir, pkgName string) error {
	var models []SegmentModel
	for _, seg := range g.schema.Segments() {
		m := SegmentModel{
			ID:          seg.ID,
			Name:        seg.Name,
			Description: CommentLines(seg.Description),
		}

		d := NewDeduper()
		for idx, f := range seg.Fields {
			pos := idx + 1
			if f.Position != "" {
				parts := strings.Split(f.Position, ".")
				if len(parts) == 2 {
					if p, err := strconv.Atoi(parts[1]); err == nil {
						pos = p
					}
				}
			}

			goName := d.Name(f.Name)
			if f.Name == "" {
				goName = d.Name(f.ID)
			}
			fType := f.DataType
			if fType == "" {
				fType = "ST"
			}
			if f.Repeatable() {
				fType = "[]" + fType
			}

			tagOpts := ""
			if f.Required() {
				tagOpts = ",required"
			}

			m.Fields = append(m.Fields, SegmentFieldModel{
				GoName: goName,
				GoType: fType,
				Tag:    fmt.Sprintf("%d%s", pos, tagOpts),
			})
		}
		models = append(models, m)
	}

	data := SegmentsFileData{
		PackageName: pkgName,
		Segments:    models,
	}

	return renderAndWrite("templates/segments.gotmpl", data, filepath.Join(targetDir, "segments.go"))
}

type groupCollector struct {
	groups []GroupModel
}

func (c *groupCollector) addGroup(msgID string, groupDeduper *Deduper, ms *schema.MessageSegment) string {
	structName := fmt.Sprintf("%s_%s", msgID, groupDeduper.Name(ms.Name))
	gm := GroupModel{
		StructName: structName,
	}

	subDeduper := NewDeduper()
	for _, child := range ms.Segments {
		var fType string
		var tagID string

		if child.Group {
			childStructName := c.addGroup(msgID, groupDeduper, child)
			fType = childStructName
			tagID = child.Name
		} else {
			fType = child.ID
			tagID = child.ID
		}

		if child.Repeatable() {
			fType = "[]" + fType
		}

		tagOpts := ""
		if child.Required() {
			tagOpts = ",required"
		}

		gm.Fields = append(gm.Fields, MessageFieldModel{
			GoName: subDeduper.Name(child.Name),
			GoType: fType,
			Tag:    fmt.Sprintf("%s%s", tagID, tagOpts),
		})
	}

	c.groups = append(c.groups, gm)
	return structName
}

func (g *Generator) GenerateMessages(targetDir, pkgName string) error {
	var models []MessageModel
	for _, msg := range g.schema.Messages() {
		m := MessageModel{
			ID:          msg.ID,
			Name:        msg.Name,
			Description: CommentLines(msg.Description),
		}

		collector := &groupCollector{}
		d := NewDeduper()
		groupDeduper := NewDeduper()

		for _, seg := range msg.Segments {
			var fType string
			var tagID string

			if seg.Group {
				groupStructName := collector.addGroup(msg.ID, groupDeduper, seg)
				fType = groupStructName
				tagID = seg.Name
			} else {
				fType = seg.ID
				tagID = seg.ID
			}

			if seg.Repeatable() {
				fType = "[]" + fType
			}

			tagOpts := ""
			if seg.Required() {
				tagOpts = ",required"
			}

			m.Fields = append(m.Fields, MessageFieldModel{
				GoName: d.Name(seg.Name),
				GoType: fType,
				Tag:    fmt.Sprintf("%s%s", tagID, tagOpts),
			})
		}

		m.Groups = collector.groups
		models = append(models, m)
	}

	data := MessagesFileData{
		PackageName: pkgName,
		Messages:    models,
	}

	return renderAndWrite("templates/messages.gotmpl", data, filepath.Join(targetDir, "messages.go"))
}

func (g *Generator) GenerateRegistry(targetDir, pkgName string) error {
	verEnum := "Version" + strings.ReplaceAll(g.ver.String(), ".", "")
	data := RegistryFileData{
		PackageName: pkgName,
		VersionEnum: verEnum,
	}

	return renderAndWrite("templates/registry.gotmpl", data, filepath.Join(targetDir, "registry.go"))
}

func (g *Generator) GenerateMessageTests(targetDir, pkgName string) error {
	data := MessageTestsFileData{
		PackageName: pkgName,
		Version:     g.ver.String(),
	}

	return renderAndWrite("templates/messages_test.gotmpl", data, filepath.Join(targetDir, "messages_test.go"))
}

func renderAndWrite(tmplPath string, data any, targetPath string) error {
	tmplContent, err := embeddedTemplates.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(tmplPath).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		_ = os.WriteFile(targetPath, buf.Bytes(), 0644)
		return fmt.Errorf("formatting error for %s: %w", targetPath, err)
	}

	return os.WriteFile(targetPath, formatted, 0644)
}
