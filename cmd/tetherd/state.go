package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"ukiran.com/tetherd/internal/procnet"
)

type State struct {
	mu             sync.Mutex
	Usage          procnet.Data
	LastSeenProc   procnet.Data
	MidNightOffset procnet.Data
	YearDay        int
	LastUpdate     time.Time
	file           *os.File
}

func NewState(path string) (*State, error) {
	sf, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	return &State{file: sf}, nil
}

func (s *State) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	currProc, err := procnet.ReadTotalDataUsage()
	if err != nil {
		return err
	}
	fmt.Println(currProc.PrintHuman()) // DEBUG:

	now := time.Now()
	today := now.YearDay()

	// load existing state, if does not exist, save zeroed values
	if err := s.persistRead(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.Usage = currProc
			s.LastSeenProc = currProc
			s.YearDay = today
			s.LastUpdate = now

			return s.persistWrite()
		}
		return err
	}

	// new day logic
	if s.YearDay != today {
		// calculate how long ago the last save was
		gap := now.Sub(s.LastUpdate)

		// If the system was off/daemon dead for more than 5 mins,
		// reset Usage
		if s.LastUpdate.IsZero() || gap > 5*time.Minute {
			s.Usage = procnet.Data{}
		} else {
			// CONTINUOUS: System was on. Add any "untracked" bytes from the
			// end of yesterday.
			delta := currProc.Delta(s.LastSeenProc)
			s.Usage = s.Usage.Add(delta)
		}
		s.YearDay = today
		s.MidNightOffset = currProc
	}

	// update bookmarks for the next Sync cycle
	s.LastSeenProc = currProc
	s.LastUpdate = now

	// save the initial state back to disk
	return s.persistWrite()
}

func (s *State) Sync(currProc procnet.Data) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	today := now.YearDay()

	if s.YearDay != today {
		s.handleDayRollover(currProc, now)
	}

	// calculate delta since last check
	var delta procnet.Data
	if currProc.GreaterOrEq(s.LastSeenProc) {
		// normal operation: OS counter grew
		delta = currProc.Delta(s.LastSeenProc)
	} else {
		// REBOOT happened: OS counter reset to zero
		// Everything in currProc is "new" since the reboot
		delta = currProc
	}

	// update state
	s.Usage = s.Usage.Add(delta)
	s.LastSeenProc = currProc
	s.LastUpdate = now

	return s.persistWrite()
}

func (s *State) handleDayRollover(currProc procnet.Data, now time.Time) {
	gap := now.Sub(s.LastUpdate)

	if gap < 5*time.Minute {
		// Keep the data that happened between 11:59 and 12:01 By doing nothing
		// here, s.Usage remains, and the next calculation in Sync() will add
		// to it.
	} else {
		// Long gap (System was off): Start fresh for the new day
		s.Usage = procnet.Data{}
	}
	s.YearDay = now.YearDay()
	s.MidNightOffset = currProc
}

// persistWrite: assumes the caller has already locked s.mu
func (s *State) persistWrite() error {
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	// write the new state
	if err := gob.NewEncoder(s.file).Encode(s); err != nil {
		return err
	}
	// force the OS to write to physical disk
	return s.file.Sync()
}

// persistRead: assumes the caller has already locked s.mu
func (s *State) persistRead() error {
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error reading state file: %w", err)
	}
	err := gob.NewDecoder(s.file).Decode(s)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error decoding state file: %w", err)
	}
	return nil
}

func (s *State) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

func (s *State) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persistWrite()
}

func (s *State) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persistRead()
}
