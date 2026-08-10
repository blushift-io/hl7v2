package main

import (
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

type config struct {
	Port    int      `yaml:"port" json:"port"`
	Actions []Action `yaml:"actions" json:"actions"`
}

func defaultConfig() *config {
	return &config{
		Port: 2525,
		Actions: []Action{
			{
				Type: ActionSave,
				Options: map[string]any{
					"output_path":       "./received_messages",
					"filename_template": defaultFilenameTmpl,
				},
			},
		},
	}
}

func readConfig(fileName string) (*config, error) {
	conf := defaultConfig()

	b, err := os.ReadFile(fileName)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

// ActionType specifies the category of action to take on incoming messages.
//go:generate go run -mod=mod github.com/dmarkham/enumer -type ActionType -yaml -json -linecomment -output=enums.go
type ActionType int

const (
	// ActionNoop performs no action on incoming messages.
	ActionNoop ActionType = iota //noop
	// ActionSave saves incoming messages to disk.
	ActionSave //save
	// ActionForward forwards incoming messages to another MLLP endpoint.
	ActionForward //forward
)

// Action defines the configuration for a single message processing action.
type Action struct {
	Type    ActionType     `yaml:"type" json:"type"`
	Options map[string]any `yaml:"options" json:"options"`
}

// SaveOptions contains configuration settings for saving messages to the filesystem.
type SaveOptions struct {
	OutputPath       string `yaml:"output-path" json:"output-path"`
	FilenameTemplate string `yaml:"filename-template" json:"filename-template"`
}

// ForwardOptions contains configuration settings for forwarding messages to a remote MLLP server.
type ForwardOptions struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}
