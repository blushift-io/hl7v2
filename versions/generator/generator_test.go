package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blushift-io/hl7v2/versions/generator"
)

func TestGenerateDatatypes(t *testing.T) {
	gen, err := generator.New("2.5")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	tmpDir := t.TempDir()
	err = gen.GenerateDatatypes(tmpDir, "v25")
	if err != nil {
		t.Fatalf("failed to generate datatypes: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "datatypes.go"))
	if err != nil {
		t.Fatalf("failed to read datatypes.go: %v", err)
	}

	code := string(content)

	// Verify CE constructor is present and returns *CE
	if !strings.Contains(code, "func NewCE() *CE {") {
		t.Errorf("datatypes.go missing func NewCE() *CE")
	}

	// Verify SetIdentifier method is present
	if !strings.Contains(code, "func (d *CE) SetIdentifier(v string) *CE {") {
		t.Errorf("datatypes.go missing func (d *CE) SetIdentifier(v string) *CE")
	}

	if !strings.Contains(code, "func (d *CE) Validate() error {") {
		t.Errorf("datatypes.go missing Validate method on CE")
	}
}

func TestGenerateSegments(t *testing.T) {
	gen, err := generator.New("2.5")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	tmpDir := t.TempDir()
	err = gen.GenerateSegments(tmpDir, "v25")
	if err != nil {
		t.Fatalf("failed to generate segments: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "segments.go"))
	if err != nil {
		t.Fatalf("failed to read segments.go: %v", err)
	}

	code := string(content)

	if !strings.Contains(code, "func NewPID() *PID {") {
		t.Errorf("segments.go missing func NewPID() *PID")
	}

	if !strings.Contains(code, "func (s *PID) SetSetIDPID(v *SI) *PID {") {
		t.Errorf("segments.go missing SetSetIDPID method")
	}

	if !strings.Contains(code, "func (s *PID) AddPatientName(v XPN) *PID {") {
		t.Errorf("segments.go missing AddPatientName method")
	}

	if !strings.Contains(code, "func NewMSH(msgType *MSG) *MSH {") {
		t.Errorf("segments.go missing func NewMSH(msgType *MSG) *MSH")
	}

	if !strings.Contains(code, "func (s *PID) Validate() error {") {
		t.Errorf("segments.go missing Validate method on PID")
	}
}

func TestGenerateMessages(t *testing.T) {
	gen, err := generator.New("2.5")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	tmpDir := t.TempDir()
	err = gen.GenerateMessages(tmpDir, "v25")
	if err != nil {
		t.Fatalf("failed to generate messages: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "messages.go"))
	if err != nil {
		t.Fatalf("failed to read messages.go: %v", err)
	}

	code := string(content)

	if !strings.Contains(code, "func NewADT_A01() *ADT_A01 {") {
		t.Errorf("messages.go missing func NewADT_A01() *ADT_A01")
	}

	if !strings.Contains(code, "MSH: NewMSH(NewMSG().SetMessageCode(\"ADT\").SetTriggerEvent(\"A01\").SetMessageStructure(\"ADT_A01\")),") {
		t.Errorf("messages.go missing expected MSH initialization in NewADT_A01")
	}

	if !strings.Contains(code, "func (m *ADT_A01) SetPID(v *PID) *ADT_A01 {") {
		t.Errorf("messages.go missing SetPID method")
	}

	if !strings.Contains(code, "func (m *ADT_A01) AddNK1(v NK1) *ADT_A01 {") {
		t.Errorf("messages.go missing AddNK1 method")
	}

	if !strings.Contains(code, "func (m *ADT_A01) Validate() error {") {
		t.Errorf("messages.go missing Validate method on ADT_A01")
	}
}

func TestGenerateTables(t *testing.T) {
	gen, err := generator.New("2.5")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	tmpDir := t.TempDir()
	err = gen.GenerateTables(tmpDir, "v25")
	if err != nil {
		t.Fatalf("failed to generate tables: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "tables.go"))
	if err != nil {
		t.Fatalf("failed to read tables.go: %v", err)
	}

	code := string(content)

	if !strings.Contains(code, "type Table0001 string") {
		t.Errorf("tables.go missing type Table0001 string")
	}

	if !strings.Contains(code, "Table0001Female Table0001 = \"F\"") {
		t.Errorf("tables.go missing Table0001Female constant")
	}

	if !strings.Contains(code, "func (t Table0001) Description() string {") {
		t.Errorf("tables.go missing Description method")
	}

	if !strings.Contains(code, "func (t Table0001) IsValid() bool {") {
		t.Errorf("tables.go missing IsValid method")
	}

	if !strings.Contains(code, "func (t Table0001) IS() *IS {") {
		t.Errorf("tables.go missing IS method")
	}

	if !strings.Contains(code, "func (t Table0001) CE() *CE {") {
		t.Errorf("tables.go missing CE method")
	}
}
