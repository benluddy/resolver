package cache

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/errors"
)

type Cache struct {
	logger                   logrus.StdLogger
	providers                []SourceProvider
	entries                  map[SourceKey]*entry
	ttl                      time.Duration
	sem                      chan struct{}
	m                        sync.RWMutex
	snapshotConcurrencyLimit int
}

type Option func(*Cache)

func WithSnapshotConcurrencyLimit(l int) Option {
	return func(c *Cache) {
		c.snapshotConcurrencyLimit = l
	}
}

func WithLogger(logger logrus.StdLogger) Option {
	return func(c *Cache) {
		c.logger = logger
	}
}

func WithSourceProviders(providers ...SourceProvider) Option {
	return func(c *Cache) {
		c.providers = providers
	}
}

func New(opts ...Option) *Cache {
	const (
		DefaultSnapshotConcurrencyLimit = 4
		DefaultSnapshotTTL              = 5 * time.Minute
	)

	c := Cache{
		logger: func() logrus.StdLogger {
			logger := logrus.New()
			logger.SetOutput(io.Discard)
			return logger
		}(),
		entries:                  make(map[SourceKey]*entry),
		ttl:                      DefaultSnapshotTTL,
		sem:                      make(chan struct{}, DefaultSnapshotConcurrencyLimit),
		snapshotConcurrencyLimit: DefaultSnapshotConcurrencyLimit,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

type OperatorFinder interface {
	Find(Predicate) []Operator
	Error() error
}

// return a snapshot collection that includes only data from sources whose labels match selector. if not cached, may begin doing work here
func (c *Cache) Fetch(selector func(map[string]string) bool) OperatorFinder {
	const (
		CachePopulateTimeout = time.Minute
	)

	now := time.Now()

	var sources []LabeledSource
	for _, provider := range c.providers {
		sources = append(sources, provider.Sources()...)

	}

	result := NamespacedOperatorCache{
		entries: make(map[SourceKey]*entry, len(sources)),
	}

	var misses []LabeledSource
	func() {
		c.m.RLock()
		defer c.m.RUnlock()
		for _, source := range sources {
			if entry, ok := c.entries[Key(source)]; ok {
				if !entry.Expired(now) {
					result.entries[Key(source)] = entry
				} else {
					misses = append(misses, source)
				}
			} else {
				misses = append(misses, source)
			}
		}
	}()

	if len(misses) == 0 {
		return &result
	}

	c.m.Lock()
	defer c.m.Unlock()

	// Take the opportunity to clear expired snapshots while holding the lock.
	var expired []SourceKey
	for key, entry := range c.entries {
		if entry.Expired(now) {
			entry.Cancel()
			expired = append(expired, key)
		}
	}
	for _, key := range expired {
		delete(c.entries, key)
	}

	// Check for any snapshots that were populated while waiting to acquire the lock.
	var found int
	for i := range misses {
		if entry, ok := c.entries[Key(misses[i])]; ok && !entry.Expired(now) {
			result.entries[Key(misses[i])] = entry
			misses[found], misses[i] = misses[i], misses[found]
			found++
		}
	}
	misses = misses[found:]

	for _, miss := range misses {
		ctx, cancel := context.WithTimeout(context.TODO(), CachePopulateTimeout)

		e := entry{
			expiry:       now.Add(c.ttl),
			pop:          cancel,
			sourceLabels: miss.Labels,
		}
		e.m.Lock()
		c.entries[Key(miss)] = &e
		result.entries[Key(miss)] = &e

		go func(ctx context.Context, e *entry, miss Source) {
			defer e.m.Unlock()

			c.sem <- struct{}{}
			defer func() { <-c.sem }()

			e.snapshot, e.err = miss.Snapshot(ctx)
		}(ctx, &e, miss)
	}

	return &result
}

type snapshots []*Snapshot

type NamespacedOperatorCache struct {
	entries map[SourceKey]*entry
}

func (c *NamespacedOperatorCache) Error() error {
	var errs []error
	for key, entry := range c.entries {
		entry.m.RLock()
		err := entry.err
		entry.m.RUnlock()
		if err != nil {
			errs = append(errs, fmt.Errorf("error populating content from source %q:: %w", key, err))
		}
	}
	return errors.NewAggregate(errs)
}

// func (c *NamespacedOperatorCache) FindPreferred(comparer SnapshotComparer, p OperatorPredicate) []Operator {
// 	snapshots := make([]*Snapshot, len(c.snapshots))
// 	for _, snapshot := range c.snapshots {
// 		snapshots = append(snapshots, snapshot.snapshot)
// 	}
// 	sort.Slice(snapshots, func(i, j int) bool {
// 		return false
// 	})
// 	var result []Operator
// 	for _, snapshot := range snapshots {
// 		result = append(result, snapshot.Find(p))
// 	}
// 	return result
// }

// can pass in a static Source at cache creation time instead of this
// func (c *NamespacedOperatorCache) WithExistingOperators(snapshot *CatalogSnapshot) MultiCatalogOperatorFinder {
// 	o := &NamespacedOperatorCache{
// 		namespaces: c.namespaces,
// 		existing:   &snapshot.key,
// 		snapshots:  c.snapshots,
// 	}
// 	o.snapshots[snapshot.key] = snapshot
// 	return o
// }

func (c *NamespacedOperatorCache) Find(p Predicate) []Operator {
	return nil
}

type entry struct {
	m            sync.RWMutex
	expiry       time.Time
	pop          context.CancelFunc
	err          error
	snapshot     *Snapshot
	sourceLabels map[string]string
}

func (s *entry) Find(p Predicate) []Operator {
	s.m.RLock()
	defer s.m.Unlock()
	return Filter(s.snapshot.Operators, p)
}

func (s *entry) Error() error {
	s.m.RLock()
	defer s.m.Unlock()
	return s.err
}

func (s *entry) Cancel() {
	s.pop()
}

func (s *entry) Expired(at time.Time) bool {
	s.m.RLock()
	defer s.m.Unlock()
	return !at.Before(s.expiry) || !s.snapshot.Valid()
}

// NewRunningOperatorSnapshot creates a CatalogSnapshot that represents a set of existing installed operators
// in the cluster.
// func NewRunningOperatorSnapshot(logger logrus.FieldLogger, key registry.CatalogKey, o []*Operator) *CatalogSnapshot {
// 	return &CatalogSnapshot{
// 		logger:    logger,
// 		key:       key,
// 		operators: o,
// 	}
// }

// type SortableSnapshots []Snapshot

// func NewSortableSnapshots(namespaces []string, snapshots map[registry.CatalogKey]*CatalogSnapshot) SortableSnapshots {
// 	sorted := SortableSnapshots{
// 		existing:   existing,
// 		preferred:  preferred,
// 		snapshots:  make([]*CatalogSnapshot, 0),
// 		namespaces: make(map[string]int, 0),
// 	}
// 	for i, n := range namespaces {
// 		sorted.namespaces[n] = i
// 	}
// 	for _, s := range snapshots {
// 		sorted.snapshots = append(sorted.snapshots, s)
// 	}
// 	return sorted
// }

// var _ sort.Interface = SortableSnapshots{}

// // Len is the number of elements in the collection.
// func (s SortableSnapshots) Len() int {
// 	return len(s.snapshots)
// }

// // Less reports whether the element with
// // index i should sort before the element with index j.
// func (s SortableSnapshots) Less(i, j int) bool {
// 	// existing operators are preferred over catalog operators
// 	if s.existing != nil &&
// 		s.snapshots[i].key.Name == s.existing.Name &&
// 		s.snapshots[i].key.Namespace == s.existing.Namespace {
// 		return true
// 	}
// 	if s.existing != nil &&
// 		s.snapshots[j].key.Name == s.existing.Name &&
// 		s.snapshots[j].key.Namespace == s.existing.Namespace {
// 		return false
// 	}

// 	// preferred catalog is less than all other catalogs
// 	if s.preferred != nil &&
// 		s.snapshots[i].key.Name == s.preferred.Name &&
// 		s.snapshots[i].key.Namespace == s.preferred.Namespace {
// 		return true
// 	}
// 	if s.preferred != nil &&
// 		s.snapshots[j].key.Name == s.preferred.Name &&
// 		s.snapshots[j].key.Namespace == s.preferred.Namespace {
// 		return false
// 	}

// 	// the rest are sorted first on priority, namespace and then by name
// 	if s.snapshots[i].priority != s.snapshots[j].priority {
// 		return s.snapshots[i].priority > s.snapshots[j].priority
// 	}
// 	if s.snapshots[i].key.Namespace != s.snapshots[j].key.Namespace {
// 		return s.namespaces[s.snapshots[i].key.Namespace] < s.namespaces[s.snapshots[j].key.Namespace]
// 	}

// 	return s.snapshots[i].key.Name < s.snapshots[j].key.Name
// }

// // Swap swaps the elements with indexes i and j.
// func (s SortableSnapshots) Swap(i, j int) {
// 	s.snapshots[i], s.snapshots[j] = s.snapshots[j], s.snapshots[i]
// }

// type MultiCatalogOperatorFinder interface {
// 	Catalog(registry.CatalogKey) OperatorFinder
// 	FindPreferred(SnapshotComparer, ...OperatorPredicate) []*Operator
// 	//	WithExistingOperators(*CatalogSnapshot) MultiCatalogOperatorFinder
// 	Error() error
// 	OperatorFinder
// }

func AtLeast(n int, operators []Operator) ([]Operator, error) {
	if len(operators) < n {
		return nil, fmt.Errorf("expected at least %d operator(s), got %d", n, len(operators))
	}
	return operators, nil
}

func ExactlyOne(operators []Operator) (*Operator, error) {
	if len(operators) != 1 {
		return nil, fmt.Errorf("expected exactly one operator, got %d", len(operators))
	}
	return &operators[0], nil
}

func Filter(operators []Operator, p Predicate) []Operator {
	var result []Operator
	for _, o := range operators {
		if p.Test(&o) {
			result = append(result, o)
		}
	}
	return result
}
