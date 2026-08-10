// Package script provides JavaScript-based transformation capabilities for HL7 v2 messages using the Goja ECMAScript engine.
package script

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/blushift-io/hl7v2"
	"github.com/dop251/goja"
)

//go:embed lib.js
var lib []byte

func newScript(userFn string) (string, error) {
	scr := strings.ReplaceAll(string(lib), "// script body", userFn)
	t, err := template.New("").Parse(scr)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %s", err.Error())
	}

	buf := bytes.NewBuffer(nil)
	vars := struct {
		ScriptBody string
	}{
		ScriptBody: userFn,
	}

	if err := t.Execute(buf, vars); err != nil {
		return "", fmt.Errorf("error executing template: %s", err.Error())
	}

	return buf.String(), nil
}

// NewTransform compiles a JavaScript transformation script into a RawTransform function.
func NewTransform(script string) (hl7v2.RawTransform, error) {
	scr, err := newScript(script)
	if err != nil {
		return nil, err
	}

	return func(msg []byte) ([]byte, error) {
		s := string(msg)

		vm := goja.New()
		if err := vm.Set("rawMessage", s); err != nil {
			return nil, fmt.Errorf("error setting rawMessage: %s", err.Error())
		}

		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

		logFn := func(call goja.FunctionCall) goja.Value {
			var args []string
			for _, arg := range call.Arguments {
				args = append(args, arg.String())
			}

			msg := strings.Join(args, " ") + "\n"
			fmt.Println(msg)

			return goja.Undefined()
		}

		if err := vm.Set("log", logFn); err != nil {
			return nil, fmt.Errorf("error setting log function: %w", err)
		}

		consoleObj := vm.NewObject()
		if err := consoleObj.Set("log", logFn); err != nil {
			return nil, fmt.Errorf("error setting log function on console object: %w", err)
		}

		if err := vm.Set("console", consoleObj); err != nil {
			return nil, fmt.Errorf("error setting console object: %w", err)
		}

		if err := vm.Set("getHeader", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				panic(vm.ToValue("getHeader requires at least one argument"))
			}
			b := []byte(call.Arguments[0].String())
			header, err := hl7v2.ParseHeader(b)
			if err != nil {
				log.Printf("error parsing message header: %v", err)
				panic(vm.ToValue(err.Error()))
			}

			res := vm.ToValue(header)
			return res
		}); err != nil {
			return nil, fmt.Errorf("error setting getHeader function: %w", err)
		}

		if err := vm.Set("parseMessage", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				panic(vm.ToValue("parseMessage requires at least one argument"))
			}

			b := []byte(call.Arguments[0].String())
			msg, err := hl7v2.ParseRaw(b)
			if err != nil {
				panic(vm.ToValue(err.Error()))
			}

			res := vm.ToValue(msg)
			return res
		}); err != nil {
			return nil, fmt.Errorf("error setting parseMessage function: %w", err)
		}

		v, err := vm.RunString(scr)
		if err != nil {
			return nil, fmt.Errorf("error running script: %s", err.Error())
		}

		res := v.ToString().String()

		return []byte(res), nil
	}, nil
}
