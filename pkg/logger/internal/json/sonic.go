//go:build sonic && (linux || windows || darwin) && amd64

package json

import "github.com/bytedance/sonic"

var NewEncoder = sonic.ConfigDefault.NewEncoder
