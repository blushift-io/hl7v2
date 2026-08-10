// Package generator provides code generation tools for creating version-specific HL7v2 Go structures from schema specifications.
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

// SpecType represents the type of HL7v2 schema specification entity.
type SpecType int

const (
	// SpecTypeUnknown indicates an unspecified schema entity type.
	SpecTypeUnknown SpecType = iota
	// SpecTypeDatatype indicates a data type schema entity.
	SpecTypeDatatype
	// SpecTypeTable indicates a table schema entity.
	SpecTypeTable
	// SpecTypeSegment indicates a segment schema entity.
	SpecTypeSegment
	// SpecTypeMessage indicates a message schema entity.
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

// Deduper tracks generated struct field and type names to avoid name collisions.
type Deduper struct {
	used map[string]int
	seen map[string]bool
}

// NewDeduper creates a new Deduper instance.
func NewDeduper() *Deduper {
	return &Deduper{
		used: make(map[string]int),
		seen: make(map[string]bool),
	}
}

// Name generates a sanitized, unique Go identifier for a raw string name.
func (d *Deduper) Name(raw string) string {
	base := SanitizeIdentifier(raw)
	candidate := base
	count := d.used[base]

	for d.seen[candidate] {
		count++
		candidate = fmt.Sprintf("%s_%d", base, count)
	}

	d.used[base] = count
	d.seen[candidate] = true
	return candidate
}

// ParamDeduper manages parameter name deduplication.
type ParamDeduper struct {
	used map[string]int
}

// NewParamDeduper creates a new ParamDeduper instance.
func NewParamDeduper() *ParamDeduper {
	return &ParamDeduper{used: make(map[string]int)}
}

// Name returns a unique parameter name by appending a numeric suffix if repeated.
func (d *ParamDeduper) Name(raw string) string {
	d.used[raw]++
	count := d.used[raw]
	if count == 1 {
		return raw
	}
	return fmt.Sprintf("%s%d", raw, count)
}

// SanitizeIdentifier cleans and formats a raw string into a valid exported Go identifier.
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

func uncapitalizeIdentifier(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "param"
	}
	runes := []rune(s)
	if len(runes) == 0 {
		return "param"
	}

	i := 0
	for i < len(runes) && unicode.IsUpper(runes[i]) {
		i++
	}

	var res string
	if i == 0 {
		res = s
	} else if i == 1 {
		res = string(unicode.ToLower(runes[0])) + string(runes[1:])
	} else if i == len(runes) {
		res = strings.ToLower(s)
	} else {
		res = strings.ToLower(string(runes[:i-1])) + string(runes[i-1:])
	}

	if goKeywords[strings.ToLower(res)] {
		res += "_"
	}

	return res
}

// CommentLines formats a multi-line string into Go comment lines prefixed with //.
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

// PackageName converts a version string like "v2.5.1" into a Go package name like "v251".
func PackageName(ver string) string {
	ver = strings.TrimPrefix(ver, "v")
	ver = strings.ReplaceAll(ver, ".", "")
	return "v" + ver
}

// Generator handles code generation for a specific HL7v2 version schema.
type Generator struct {
	ver    hl7v2.Version
	schema *schema.Schema
	opts   *Options
}

// New creates a new Generator for the specified HL7v2 version string and options.
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

// DatatypeFieldModel holds template data for a field within a data type.
type DatatypeFieldModel struct {
	GoName      string
	GoType      string
	Tag         string
	ParamType   string
	InitExpr    string
	MethodName  string
	IsRequired  bool
	IsPrimitive bool
	IsSlice     bool
}

// DatatypeModel holds template data for an HL7v2 data type definition.
type DatatypeModel struct {
	ID          string
	Name        string
	Description string
	IsPrimitive bool
	Fields      []DatatypeFieldModel
}

// DatatypesFileData holds template rendering data for a datatypes.go file.
type DatatypesFileData struct {
	PackageName         string
	Datatypes           []DatatypeModel
	HasComplexDatatypes bool
}

// SegmentFieldModel holds template data for a field within a segment.
type SegmentFieldModel struct {
	GoName      string
	GoType      string
	Tag         string
	IsSlice     bool
	IsPrimitive bool
	IsRequired  bool
	ElemType    string
	MethodName  string
}

// MSHInitModel holds template data for MSH segment initialization logic.
type MSHInitModel struct {
	FieldSeparatorName     string
	EncodingCharactersName string
	DateTimeName           string
	DateTimeExpr           string
	MessageTypeName        string
	MessageTypeParamType   string
	MessageControlIDName   string
	ProcessingIDName       string
	ProcessingIDExpr       string
	VersionIDName          string
	VersionIDExpr          string
}

// SegmentModel holds template data for an HL7v2 segment definition.
type SegmentModel struct {
	ID          string
	Name        string
	Description string
	Fields      []SegmentFieldModel
	IsMSH       bool
	MSHInit     MSHInitModel
}

// SegmentsFileData holds template rendering data for a segments.go file.
type SegmentsFileData struct {
	PackageName string
	Segments    []SegmentModel
}

// MessageFieldModel holds template data for a field within a message.
type MessageFieldModel struct {
	GoName     string
	GoType     string
	Tag        string
	IsSlice    bool
	IsRequired bool
	ElemType   string
	MethodName string
}

// GroupModel holds template data for a message group structure.
type GroupModel struct {
	StructName string
	Fields     []MessageFieldModel
}

// MessageModel holds template data for an HL7v2 message definition.
type MessageModel struct {
	ID          string
	Name        string
	Description string
	Groups      []GroupModel
	Fields      []MessageFieldModel
	MSHInitExpr string
}

// MessagesFileData holds template rendering data for a messages.go file.
type MessagesFileData struct {
	PackageName string
	Messages    []MessageModel
}

// RegistryFileData holds template rendering data for a registry.go file.
type RegistryFileData struct {
	PackageName string
	VersionEnum string
}

// MessageTestsFileData holds template rendering data for a messages_test.go file.
type MessageTestsFileData struct {
	PackageName string
	Version     string
}

// TableEntryModel holds template data for an entry within a table.
type TableEntryModel struct {
	ConstName          string
	Value              string
	Description        string
	DescriptionEscaped string
	Comment            string
}

// TableModel holds template data for an HL7v2 table definition.
type TableModel struct {
	ID          string
	TypeName    string
	Name        string
	Description string
	Entries     []TableEntryModel
}

// TablesFileData holds template rendering data for a tables.go file.
type TablesFileData struct {
	PackageName string
	Tables      []TableModel
	HasIS       bool
	HasID       bool
	HasST       bool
	HasCE       bool
	HasCWE      bool
	HasCNE      bool
}

// GenerateAll generates all Go source files for the configured version under outputBase.
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
	if err := g.GenerateTables(targetDir, pkgName); err != nil {
		return fmt.Errorf("tables gen error: %w", err)
	}
	if err := g.GenerateRegistry(targetDir, pkgName); err != nil {
		return fmt.Errorf("registry gen error: %w", err)
	}
	if err := g.GenerateMessageTests(targetDir, pkgName); err != nil {
		return fmt.Errorf("message tests gen error: %w", err)
	}

	return nil
}

// GenerateDatatypes generates the datatypes.go file for a target directory and package.
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
				methodName := "Set" + goName

				fType := f.DataType
				if fType == "" {
					fType = "ST"
				}

				fDT := g.schema.DataType(fType)
				isPrimitive := fDT == nil || fDT.IsPrimitive()

				var goType string
				var paramType string
				var initExpr string

				if isPrimitive {
					if f.Repeatable() {
						goType = "[]" + fType
						paramType = "[]string"
						initExpr = fmt.Sprintf("func() []%s { if v == nil { return nil }; res := make([]%s, len(v)); for i, val := range v { res[i] = %s(val) }; return res }()", fType, fType, fType)
					} else {
						goType = fType
						paramType = "string"
						initExpr = fmt.Sprintf("%s(v)", fType)
					}
				} else {
					if f.Repeatable() {
						goType = "[]" + fType
						paramType = "[]" + fType
						initExpr = "v"
					} else {
						goType = fType
						paramType = fType
						initExpr = "v"
					}
				}

				m.Fields = append(m.Fields, DatatypeFieldModel{
					GoName:      goName,
					GoType:      goType,
					Tag:         fmt.Sprintf("%d", pos),
					ParamType:   paramType,
					InitExpr:    initExpr,
					MethodName:  methodName,
					IsRequired:  f.Required(),
					IsPrimitive: isPrimitive,
					IsSlice:     f.Repeatable(),
				})
			}
		}
		models = append(models, m)
	}

	hasComplex := false
	for _, m := range models {
		if !m.IsPrimitive {
			for _, f := range m.Fields {
				if f.IsRequired || f.IsSlice || !f.IsPrimitive {
					hasComplex = true
					break
				}
			}
			if hasComplex {
				break
			}
		}
	}

	data := DatatypesFileData{
		PackageName:         pkgName,
		Datatypes:           models,
		HasComplexDatatypes: hasComplex,
	}

	return renderAndWrite("templates/datatypes.gotmpl", data, filepath.Join(targetDir, "datatypes.go"))
}

// GenerateSegments generates the segments.go file for a target directory and package.
func (g *Generator) GenerateSegments(targetDir, pkgName string) error {
	var models []SegmentModel
	for _, seg := range g.schema.Segments() {
		m := SegmentModel{
			ID:          seg.ID,
			Name:        seg.Name,
			Description: CommentLines(seg.Description),
			IsMSH:       seg.ID == "MSH",
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

			isSlice := f.Repeatable()
			elemType := fType
			goType := fType
			methodName := "Set" + goName
			if isSlice {
				goType = "[]" + fType
				methodName = "Add" + goName
			}

			tagOpts := ""
			if f.Required() {
				tagOpts = ",required"
			}

			fDT := g.schema.DataType(fType)
			isPrimitive := fDT == nil || fDT.IsPrimitive()

			m.Fields = append(m.Fields, SegmentFieldModel{
				GoName:      goName,
				GoType:      goType,
				Tag:         fmt.Sprintf("%d%s", pos, tagOpts),
				IsSlice:     isSlice,
				IsPrimitive: isPrimitive,
				IsRequired:  f.Required(),
				ElemType:    elemType,
				MethodName:  methodName,
			})
		}

		if m.IsMSH {
			init := MSHInitModel{}

			for _, f := range m.Fields {
				if strings.HasPrefix(f.GoName, "FieldSeparator") {
					init.FieldSeparatorName = f.GoName
				} else if strings.HasPrefix(f.GoName, "EncodingCharacters") {
					init.EncodingCharactersName = f.GoName
				} else if strings.HasPrefix(f.GoName, "DateTime") {
					init.DateTimeName = f.GoName
					fDT := g.schema.DataType(f.ElemType)
					if fDT == nil || fDT.IsPrimitive() {
						init.DateTimeExpr = fmt.Sprintf("New%s(time.Now().Format(\"20060102150405\"))", f.ElemType)
					} else {
						subF := fDT.Fields[0]
						subGoName := SanitizeIdentifier(subF.Name)
						if subF.Name == "" {
							subGoName = SanitizeIdentifier(subF.ID)
						}
						init.DateTimeExpr = fmt.Sprintf("New%s().Set%s(time.Now().Format(\"20060102150405\"))", f.ElemType, subGoName)
					}
				} else if strings.HasPrefix(f.GoName, "MessageType") {
					init.MessageTypeName = f.GoName
					init.MessageTypeParamType = f.ElemType
				} else if strings.HasPrefix(f.GoName, "MessageControl") {
					init.MessageControlIDName = f.GoName
				} else if strings.HasPrefix(f.GoName, "ProcessingI") {
					init.ProcessingIDName = f.GoName
					fDT := g.schema.DataType(f.ElemType)
					if fDT == nil || fDT.IsPrimitive() {
						init.ProcessingIDExpr = fmt.Sprintf("New%s(\"P\")", f.ElemType)
					} else if len(fDT.Fields) > 0 {
						subF := fDT.Fields[0]
						subGoName := SanitizeIdentifier(subF.Name)
						if subF.Name == "" {
							subGoName = SanitizeIdentifier(subF.ID)
						}
						init.ProcessingIDExpr = fmt.Sprintf("New%s().Set%s(\"P\")", f.ElemType, subGoName)
					}
				} else if strings.HasPrefix(f.GoName, "VersionI") {
					init.VersionIDName = f.GoName
					fDT := g.schema.DataType(f.ElemType)
					if fDT == nil || fDT.IsPrimitive() {
						init.VersionIDExpr = fmt.Sprintf("New%s(\"%s\")", f.ElemType, g.ver.String())
					} else {
						subF := fDT.Fields[0]
						subGoName := SanitizeIdentifier(subF.Name)
						if subF.Name == "" {
							subGoName = SanitizeIdentifier(subF.ID)
						}
						init.VersionIDExpr = fmt.Sprintf("New%s().Set%s(\"%s\")", f.ElemType, subGoName, g.ver.String())
					}
				}
			}

			if init.FieldSeparatorName == "" {
				init.FieldSeparatorName = "FieldSeparator"
			}
			if init.EncodingCharactersName == "" {
				init.EncodingCharactersName = "EncodingCharacters"
			}
			if init.DateTimeName == "" {
				init.DateTimeName = "DateTimeOfMessage"
			}
			if init.DateTimeExpr == "" {
				init.DateTimeExpr = "TS(time.Now().Format(\"20060102150405\"))"
			}
			if init.MessageTypeName == "" {
				init.MessageTypeName = "MessageType"
			}
			if init.MessageTypeParamType == "" {
				init.MessageTypeParamType = "MSG"
			}
			if init.MessageControlIDName == "" {
				init.MessageControlIDName = "MessageControlID"
			}
			if init.VersionIDName == "" {
				init.VersionIDName = "VersionID"
			}
			if init.VersionIDExpr == "" {
				init.VersionIDExpr = fmt.Sprintf("VID{VersionID: ID(\"%s\")}", g.ver.String())
			}

			m.MSHInit = init
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

		goName := subDeduper.Name(child.Name)
		isSlice := child.Repeatable()
		elemType := fType
		goType := fType
		methodName := "Set" + goName
		if isSlice {
			goType = "[]" + fType
			methodName = "Add" + goName
		}

		tagOpts := ""
		if child.Required() {
			tagOpts = ",required"
		}

		gm.Fields = append(gm.Fields, MessageFieldModel{
			GoName:     goName,
			GoType:     goType,
			Tag:        fmt.Sprintf("%s%s", tagID, tagOpts),
			IsSlice:    isSlice,
			IsRequired: child.Required(),
			ElemType:   elemType,
			MethodName: methodName,
		})
	}

	c.groups = append(c.groups, gm)
	return structName
}

// GenerateMessages generates the messages.go file for a target directory and package.
func (g *Generator) GenerateMessages(targetDir, pkgName string) error {
	var mshMsgType string
	mshSeg := g.schema.Segment("MSH")
	if mshSeg != nil {
		for _, f := range mshSeg.Fields {
			if strings.HasPrefix(f.Name, "Message Type") || f.ID == "MSG" || f.ID == "CM_MSG" || f.ID == "ID" {
				mshMsgType = f.DataType
				break
			}
		}
	}
	if mshMsgType == "" {
		mshMsgType = "MSG"
	}

	var models []MessageModel
	for _, msg := range g.schema.Messages() {
		parts := strings.Split(msg.ID, "_")
		code := parts[0]
		event := ""
		if len(parts) > 1 {
			event = parts[1]
		}
		structure := msg.ID

		var mshInitExpr string
		switch mshMsgType {
		case "MSG":
			f1Setter := "SetMessageCode"
			msgDT := g.schema.DataType("MSG")
			if msgDT != nil && len(msgDT.Fields) > 0 {
				d := NewDeduper()
				f1Name := d.Name(msgDT.Fields[0].Name)
				if msgDT.Fields[0].Name == "" {
					f1Name = d.Name(msgDT.Fields[0].ID)
				}
				f1Setter = "Set" + f1Name
			}
			mshInitExpr = fmt.Sprintf("NewMSH(NewMSG().%s(\"%s\").SetTriggerEvent(\"%s\").SetMessageStructure(\"%s\"))", f1Setter, code, event, structure)
		case "CM_MSG":
			f1Setter := "SetMessageType"
			cmMsgDT := g.schema.DataType("CM_MSG")
			if cmMsgDT != nil && len(cmMsgDT.Fields) > 0 {
				d := NewDeduper()
				f1Name := d.Name(cmMsgDT.Fields[0].Name)
				if cmMsgDT.Fields[0].Name == "" {
					f1Name = d.Name(cmMsgDT.Fields[0].ID)
				}
				f1Setter = "Set" + f1Name
			}
			mshInitExpr = fmt.Sprintf("NewMSH(NewCM_MSG().%s(\"%s\").SetTriggerEvent(\"%s\"))", f1Setter, code, event)
		default:
			mshInitExpr = fmt.Sprintf("NewMSH(NewID(\"%s\"))", msg.ID)
		}

		m := MessageModel{
			ID:          msg.ID,
			Name:        msg.Name,
			Description: CommentLines(msg.Description),
			MSHInitExpr: mshInitExpr,
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

			goName := d.Name(seg.Name)
			if seg.Name == "" {
				goName = d.Name(seg.ID)
			}
			isSlice := seg.Repeatable()
			elemType := fType
			goType := fType
			methodName := "Set" + goName
			if isSlice {
				goType = "[]" + fType
				methodName = "Add" + goName
			}

			tagOpts := ""
			if seg.Required() {
				tagOpts = ",required"
			}

			m.Fields = append(m.Fields, MessageFieldModel{
				GoName:     goName,
				GoType:     goType,
				Tag:        fmt.Sprintf("%s%s", tagID, tagOpts),
				IsSlice:    isSlice,
				IsRequired: seg.Required(),
				ElemType:   elemType,
				MethodName: methodName,
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

// GenerateRegistry generates the registry.go file for a target directory and package.
func (g *Generator) GenerateRegistry(targetDir, pkgName string) error {
	verEnum := "Version" + strings.ReplaceAll(g.ver.String(), ".", "")
	data := RegistryFileData{
		PackageName: pkgName,
		VersionEnum: verEnum,
	}

	return renderAndWrite("templates/registry.gotmpl", data, filepath.Join(targetDir, "registry.go"))
}

// GenerateTables generates the tables.go file for a target directory and package.
func (g *Generator) GenerateTables(targetDir, pkgName string) error {
	var models []TableModel
	tableDeduper := NewDeduper()

	for _, tbl := range g.schema.Tables() {
		rawID := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, tbl.ID)
		if rawID == "" {
			rawID = "Unknown"
		}
		typeName := tableDeduper.Name("Table" + rawID)

		tm := TableModel{
			ID:          tbl.ID,
			TypeName:    typeName,
			Name:        tbl.Name,
			Description: CommentLines(tbl.Name),
		}

		entryDeduper := NewDeduper()
		seenValues := make(map[string]bool)

		for _, e := range tbl.Entries {
			if seenValues[e.Value] {
				continue
			}
			seenValues[e.Value] = true

			var base string
			if e.Description != "" {
				base = SanitizeIdentifier(e.Description)
			}
			if base == "" {
				base = SanitizeIdentifier(e.Value)
			}
			if base == "" {
				base = "Value"
			}
			if len(base) > 0 && base[0] >= '0' && base[0] <= '9' {
				base = "V" + base
			}

			constName := typeName + entryDeduper.Name(base)
			descEscaped := strings.ReplaceAll(strings.ReplaceAll(e.Description, "\\", "\\\\"), "\"", "\\\"")

			tm.Entries = append(tm.Entries, TableEntryModel{
				ConstName:          constName,
				Value:              e.Value,
				Description:        e.Description,
				DescriptionEscaped: descEscaped,
				Comment:            e.Comment,
			})
		}

		models = append(models, tm)
	}

	data := TablesFileData{
		PackageName: pkgName,
		Tables:      models,
		HasIS:       g.schema.DataType("IS") != nil,
		HasID:       g.schema.DataType("ID") != nil,
		HasST:       g.schema.DataType("ST") != nil,
		HasCE:       g.schema.DataType("CE") != nil,
		HasCWE:      g.schema.DataType("CWE") != nil,
		HasCNE:      g.schema.DataType("CNE") != nil,
	}

	return renderAndWrite("templates/tables.gotmpl", data, filepath.Join(targetDir, "tables.go"))
}

// GenerateMessageTests generates the messages_test.go file for a target directory and package.
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
