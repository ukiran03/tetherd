package main

import (
	"encoding/gob"
	"io"

	"ukiran.com/tetherd/internal/procnet"
)

const stateFile = "/var/lib/tetherd/state.bin"

type State struct {
	file io.ReadWriter
}

// Write writes a slice of Stats to a binary file
func (s *State) Write(stats []procnet.Stats) error {
	encoder := gob.NewEncoder(s.file)
	return encoder.Encode(stats)
}

// Read reads the binary file back into a slice
func (s *State) Read() ([]procnet.Stats, error) {
	var stats []procnet.Stats
	decoder := gob.NewDecoder(s.file)
	err := decoder.Decode(&stats)
	return stats, err
}
