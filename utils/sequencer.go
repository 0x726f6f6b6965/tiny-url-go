package utils

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type Sequencer interface {
	Next() (*big.Int, error)
}

const (
	maxEpoch    uint64 = 1<<41 - 1
	maxNode     int64  = 1<<8 - 1
	maxSequence int64  = 1<<14 - 1
	shiftEpoch  uint8  = 22
	shiftNode   uint8  = 14
)

var (
	// onceInitSeq guarantee initialize only once
	onceInitSeq sync.Once
	// rootSequencer - the root sequencer
	rootSequencer *sequencer

	// ErrInvalidNode - invalid node id
	ErrInvalidNode = fmt.Errorf("invalid node id; must be 0 ≤ id < %d", maxNode)

	// ErrStartZero - the error of starting time is zero
	ErrStartZero = errors.New("the start time cannot be a zero value")

	// ErrStartFuture - the error of starting time is in the future
	ErrStartFuture = errors.New("the start time cannot be greater than the current millisecond")

	// ErrStartExceed - the start time is more than 69 years ago.
	ErrStartExceed = errors.New("the maximum life cycle of the snowflake algorithm is 69 years")
)

type sequencer struct {
	sync.Mutex
	// nodeID is the node ID that the Snowflake generator will use for the next 8 bits
	nodeID int64
	// sequence is the last 14 bits.
	sequence int64
	// baseEpoch is the start time.
	baseEpoch int64
	// currentEpoch is the current time.
	currentEpoch int64
}

func NewSequencer(nodeID int64, start time.Time) (Sequencer, error) {
	if nodeID > maxNode {
		return nil, ErrInvalidNode
	}
	start = start.UTC()

	if start.IsZero() {
		return nil, ErrStartZero
	}

	if start.After(time.Now().UTC()) {
		return nil, ErrStartFuture
	}

	if uint64(time.Now().UTC().UnixMilli()-start.UnixMilli()) > maxEpoch {
		return nil, ErrStartExceed
	}
	onceInitSeq.Do(func() {
		rootSequencer = &sequencer{
			nodeID:       nodeID,
			sequence:     0,
			baseEpoch:    start.UnixMilli(),
			currentEpoch: start.UnixMilli(),
		}
	})
	return rootSequencer, nil
}
func (s *sequencer) Next() (*big.Int, error) {
	s.Lock()
	defer s.Unlock()
	current := time.Now().UTC().UnixMilli()
	if uint64(current-s.baseEpoch) > maxEpoch {
		return nil, ErrStartExceed
	}

	if current != s.currentEpoch {
		s.sequence = 0
		s.currentEpoch = current
	} else {
		s.sequence += 1
		if s.sequence > maxSequence {
			s.sequence = 0
			s.currentEpoch += 1
		}
	}

	result := (s.currentEpoch-s.baseEpoch)<<shiftEpoch | s.nodeID<<shiftNode | s.sequence
	num := big.NewInt(result)
	return num, nil
}
