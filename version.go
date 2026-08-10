package hl7v2

import "strings"

// Version represents an HL7 v2 standard version enum.
//
//go:generate enumer -type=Version --linecomment -output=version_enums.go
type Version int

const (
	// VersionUnknown represents an unknown or unsupported HL7 version.
	VersionUnknown Version = iota //unknown
	// Version21 represents HL7 version 2.1.
	Version21                     //2.1
	// Version22 represents HL7 version 2.2.
	Version22                     //2.2
	// Version23 represents HL7 version 2.3.
	Version23                     //2.3
	// Version231 represents HL7 version 2.3.1.
	Version231                    //2.3.1
	// Version24 represents HL7 version 2.4.
	Version24                     //2.4
	// Version25 represents HL7 version 2.5.
	Version25                     //2.5
	// Version251 represents HL7 version 2.5.1.
	Version251                    //2.5.1
	// Version26 represents HL7 version 2.6.
	Version26                     //2.6
	// Version27 represents HL7 version 2.7.
	Version27                     //2.7
	// Version271 represents HL7 version 2.7.1.
	Version271                    //2.7.1
	// Version28 represents HL7 version 2.8.
	Version28                     //2.8
	// Version281 represents HL7 version 2.8.1.
	Version281                    //2.8.1
	// Version282 represents HL7 version 2.8.2.
	Version282                    //2.8.2
	// Version29 represents HL7 version 2.9.
	Version29                     //2.9
)

// GetVersion parses a version string (e.g. "2.5" or "v2.5") into a Version enum value.
func GetVersion(ver string) Version {
	ver = strings.TrimSpace(ver)
	ver = strings.TrimPrefix(ver, "v")

	v, err := VersionString(ver)
	if err != nil {
		return VersionUnknown
	}

	return v
}
