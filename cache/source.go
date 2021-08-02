package cache

import (
	"context"
	"encoding/json"
	"hash/fnv"
)

// Snapshots should be considered read-only once populated and
// returned from a Source, and may be in concurrent use.
type Snapshot struct {
	Operators []Operator
	valid     <-chan struct{}
}

func (s *Snapshot) Valid() bool {
	if s == nil {
		return false
	}

	select {
	case <-s.valid:
		return true
	default:
		return s.valid == nil
	}
}

type SnapshotBuilder struct {
	s Snapshot
}

func (b *SnapshotBuilder) Build() *Snapshot {
	s := b.s
	b.s = Snapshot{}
	return &s
}

func (b *SnapshotBuilder) Operator(o Operator) {
	b.s.Operators = append(b.s.Operators, o)
}

// snapshot creators can cancel their snapshots directly / without knowing about the cache
func (b *SnapshotBuilder) ValidChannel(c <-chan struct{}) {
	b.s.valid = c
}

type Source interface {
	Snapshot(context.Context) (*Snapshot, error)
}

type StaticSource struct {
	Operators []Operator
	Error     error
}

func (s StaticSource) Snapshot(context.Context) (*Snapshot, error) {
	return &Snapshot{Operators: s.Operators}, s.Error
}

type LabeledSource struct {
	Source
	Labels map[string]string
}

type SourceKey uint64

func Key(s LabeledSource) SourceKey {
	fnv := fnv.New64()
	if err := json.NewEncoder(fnv).Encode(s.Labels); err != nil {
		// This implementation of io.Writer never returns a
		// non-nil error, so this should be impossible.
		panic(err)
	}
	return SourceKey(fnv.Sum64())
}

type SourceSelector interface {
	Matches(labels map[string]string) bool
}

type SourceProvider interface {
	Sources(selector SourceSelector) []LabeledSource
}

type StaticSourceProvider []LabeledSource

func (p StaticSourceProvider) Sources() []LabeledSource {
	return p
}

type SourceOrder interface {
	Less(a, b LabeledSource, delegate SourceOrder) bool
}
