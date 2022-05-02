//go:build !go_json && !(sonic && (linux || windows || darwin) && amd64)

package json

import "encoding/json"

var NewEncoder = json.NewEncoder
