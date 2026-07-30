package hl7v2

import "strings"

//go:generate enumer -type=Version --linecomment -output=version_enums.go
type Version int

const (
	VersionUnknown Version = iota //unknown
	Version21                     //2.1
	Version22                     //2.2
	Version23                     //2.3
	Version231                    //2.3.1
	Version24                     //2.4
	Version25                     //2.5
	Version251                    //2.5.1
	Version26                     //2.6
	Version27                     //2.7
	Version271                    //2.7.1
	Version28                     //2.8
	Version281                    //2.8.1
	Version282                    //2.8.2
	Version29                     //2.9
)

func GetVersion(ver string) Version {
	ver = strings.TrimSpace(ver)
	ver = strings.TrimPrefix(ver, "v")

	v, err := VersionString(ver)
	if err != nil {
		return VersionUnknown
	}

	return v
}
