# Transformers & JS Scripting

`hl7v2` provides flexible transformation capabilities for sanitizing, rewriting, and transforming HL7 v2 messages. This includes built-in functions for line ending normalization and unescaping, as well as an embedded JavaScript engine for executing custom transformation scripts.

## Core Transformer Types

`hl7v2` defines two main transformer signatures:

- **`hl7v2.RawTransform`**: `func([]byte) ([]byte, error)`  
  Operates directly on raw byte streams before or during parsing.
- **`hl7v2.ValueTransform`**: `func(Element) (Value, error)`  
  Operates on individual elements (segments, fields, or components).

### Built-in Raw & Value Transformers

```go
package main

import (
	"fmt"
	"github.com/blushift-io/hl7v2"
)

func main() {
	raw := []byte("MSH|^~\\&|SEND|FAC\r\nPID|1||123\n")

	// Normalize Windows (\r\n) or Unix (\n) line endings to HL7 standard (\r)
	fixed := hl7v2.ReplaceLineEndings(raw)
	fmt.Printf("Normalized bytes: %q\n", string(fixed))
}
```

---

## JavaScript Script Transformer (`transform/script`)

For complex or dynamic message rewriting rules (e.g. routing logic, payload sanitization, or field mapping), `hl7v2` includes an embedded JavaScript runtime powered by Goja in the `transform/script` package.

`script.NewTransform(userScript)` compiles a JavaScript snippet into an `hl7v2.RawTransform` function.

### Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/transform/script"
)

var jsCode = `
function handleMessage(m) { 
	var msh = m.getSegment("MSH");
	if (msh) {
		var evt = msh.getComponent(9, 1);
		if (evt == "ADT") {
			msh.setField(12, "2.5"); // Update version ID
			msh.addField("EXTRA_HEADER");
		}
	}

	// Add custom z-segment
	m.addSegment("ZZX", "1", "CUSTOM_DATA", "ACTIVE");

	return m.toString();
}

var p = new parser();
var m = p.parse(rawMessage);
handleMessage(m);
`

func main() {
	// Compile the JS script into a RawTransform
	transformFn, err := script.NewTransform(jsCode)
	if err != nil {
		log.Fatalf("Script compile error: %v", err)
	}

	input := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260810||ADT^A01|101|P|2.1\rPID|1||12345")
	output, err := transformFn(input)
	if err != nil {
		log.Fatalf("Transform execution error: %v", err)
	}

	fmt.Println("Transformed HL7 Message:")
	fmt.Println(string(output))
}
```

---

## JavaScript API Reference

Inside the script execution environment, the following globals, classes, and helper functions are injected:

### Environment Globals
- **`rawMessage`** (`string`): The raw incoming HL7 v2 message string passed to the transformer.
- **`console.log(...)` / `log(...)`**: Logs debug output to standard output (`os.Stdout`).

### `parser` & `message` API
- **`new parser()`**: Instantiates the embedded JavaScript HL7 message parser.
- **`p.parse(rawMessageStr)`**: Parses the raw message string into a JS `message` object.

#### `message` Object Methods
- **`m.getSegment(name)`**: Retrieves the first segment matching `name` (e.g. `"MSH"`, `"PID"`).
- **`m.getSegments(name)`**: Retrieves an array of matching segments.
- **`m.addSegment(name, val1, val2, ...)`**: Appends a new segment with the specified field values.
- **`m.toString()`**: Serializes the transformed message back to an HL7 ER7 formatted string (`\r` separated).

#### `segment` Object Methods
- **`seg.getComponent(fieldIndex, componentIndex)`**: Returns string value at specified 1-based field and component index.
- **`seg.setField(fieldIndex, value)`**: Sets the field value at 1-based `fieldIndex`.
- **`seg.addField(value)`**: Appends a new field to the end of the segment.

---

## Helper Functions

### `parseMessage(rawString)`
Parses a raw message string into an `hl7v2.RawMessage` object with low-level query support:

```javascript
try {
	var rm = parseMessage(m.toString());
	var msgType = rm.queryValue("MSH.9.1");
	console.log("Parsed Message Type: " + msgType.string());
} catch (e) {
	console.log("Error querying message: " + e);
}
```

### `getHeader(rawString)`
Parses and returns the `MessageHeader` struct for quick header field extraction:

```javascript
var hdr = getHeader(rawMessage);
console.log("Sending App: " + hdr.sendingApplication);
```
