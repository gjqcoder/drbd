//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package flex

import (
	"encoding/json"
)

type exitCode int

const (
	ExitSuccess exitCode = iota
	ExitFailure
)

type exitStatus string

const (
	StatusSuccess      exitStatus = "Success"
	StatusFailure      exitStatus = "Failure"
	StatusNotSupported exitStatus = "Not supported"
)

type callEcho struct {
	Status  exitStatus `json:"status"`
	Message string     `json:"message"`
}

type initCallCapabilities struct {
	Attach bool `json:"attach"`
}

type initCallEcho struct {
	callEcho
	Capabilities initCallCapabilities `json:"capabilities"`
}

var errorArgs = callEcho{
	Status:  StatusFailure,
	Message: "Arguments number or malformated args error.",
}

var errorNotSupported = callEcho{
	Status:  StatusNotSupported,
	Message: "Operation not supported yet.",
}

// More FlexVolume specification default options
//  https://github.com/kubernetes/community/contributors/devel/flexvolume.md
type Options struct {
	FsType  string `json:"kubernetes.io/fsType"`
	ResName string `json:"resource"`
}

func parseOptions(rawOpts string) (Options, error) {
	opts := Options{}
	err := json.Unmarshal([]byte(rawOpts), &opts)
	return opts, err
}
