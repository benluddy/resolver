package cache

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/errors"

	"github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"
)

type Source interface {
	Snapshot(context.Context) (*Snapshot, error)
}

type NamedSource struct {
	Name string
	Source
}

type SnapshotProvider interface {
	Namespaced(namespaces ...string) MultiCatalogOperatorFinder
	Expire(catalog registry.CatalogKey)
}

type Cache struct {
	logger    logrus.StdLogger
	sources   []NamedSource
	snapshots map[string]*snapshotHeader
	ttl       time.Duration
	sem       chan struct{}
	m         sync.RWMutex
}

const defaultCatalogSourcePriority int = 0

type catalogSourcePriority int

type Option func(*Cache)

func WithLogger(logger logrus.StdLogger) Option {
	return func(c *Cache) {
		c.logger = logger
	}
}

func New(sources []NamedSource, opts ...Option) *Cache {
	const (
		MaxConcurrentSnapshotUpdates = 4
	)

	c := Cache{
		logger: func() logrus.StdLogger {
			logger := logrus.New()
			logger.SetOutput(io.Discard)
			return logger
		}(),
		sources:   sources,
		snapshots: make(map[string]*snapshotHeader),
		ttl:       5 * time.Minute,
		sem:       make(chan struct{}, MaxConcurrentSnapshotUpdates),
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

type NamespacedOperatorCache struct {
	namespaces []string
	existing   *registry.CatalogKey
	snapshots  map[string]*snapshotHeader
}

func (c *NamespacedOperatorCache) Error() error {
	var errs []error
	for key, snapshot := range c.snapshots {
		snapshot.m.RLock()
		err := snapshot.err
		snapshot.m.RUnlock()
		if err != nil {
			errs = append(errs, fmt.Errorf("error populating content from source %q:: %w", key, err))
		}
	}
	return errors.NewAggregate(errs)
}

func (c *Cache) Expire(sourceName string) {
	c.m.Lock()
	defer c.m.Unlock()
	s, ok := c.snapshots[sourceName]
	if !ok {
		return
	}
	s.expiry = time.Unix(0, 0)
}

func (c *Cache) Namespaced(namespaces ...string) MultiCatalogOperatorFinder {
	const (
		CachePopulateTimeout = time.Minute
	)

	now := time.Now()

	result := NamespacedOperatorCache{
		namespaces: namespaces,
		snapshots:  make(map[string]*snapshotHeader),
	}

	var misses []NamedSource
	func() {
		c.m.RLock()
		defer c.m.RUnlock()
		for _, source := range c.sources {
			snapshot, ok := c.snapshots[source.Name]
			if ok {
				func() {
					snapshot.m.RLock()
					defer snapshot.m.RUnlock()
					if !snapshot.Expired(now) && snapshot.snapshot != nil {
						result.snapshots[source.Name] = snapshot
					} else {
						misses = append(misses, source)
					}
				}()
			}
			if !ok {
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
	var expired []string
	for name, snapshot := range c.snapshots {
		if snapshot.Expired(now) {
			snapshot.Cancel()
			expired = append(expired, name)
		}
	}
	for _, name := range expired {
		delete(c.snapshots, name)
	}

	// Check for any snapshots that were populated while waiting to acquire the lock.
	var found int
	for i := range misses {
		if snapshot, ok := c.snapshots[misses[i].Name]; ok && !snapshot.Expired(now) && snapshot.snapshot != nil {
			result.snapshots[misses[i].Name] = snapshot
			misses[found], misses[i] = misses[i], misses[found]
			found++
		}
	}
	misses = misses[found:]

	for _, miss := range misses {
		ctx, cancel := context.WithTimeout(context.Background(), CachePopulateTimeout)

		s := snapshotHeader{
			key:    miss.Name,
			expiry: now.Add(c.ttl),
			pop:    cancel,
		}
		s.m.Lock()
		c.snapshots[miss.Name] = &s
		result.snapshots[miss.Name] = &s
		go c.populate(ctx, &s, miss)
	}

	return &result
}

func (c *Cache) populate(ctx context.Context, hdr *snapshotHeader, source Source) {
	defer hdr.m.Unlock()
	defer func() {
		// Don't cache an errorred snapshot.
		if hdr.err != nil {
			hdr.expiry = time.Time{}
		}
	}()

	c.sem <- struct{}{}
	defer func() { <-c.sem }()

	if snapshot, err := source.Snapshot(ctx); err != nil {
		hdr.err = err
	} else {
		hdr.snapshot = snapshot
	}
}

func (c *NamespacedOperatorCache) Catalog(k registry.CatalogKey) OperatorFinder {
	// all catalogs match the empty catalog
	if k.Empty() {
		return c
	}
	if header, ok := c.snapshots[k]; ok {
		return header.snapshot
	}
	return EmptyOperatorFinder{}
}

type ComparisonResult int

const (
	LessThan ComparisonResult = iota
	EqualTo
	GreaterThan
	Defer
)

type SnapshotComparer interface {
	Compare(a, b *Snapshot) ComparisonResult
}

func Less(a, b *Snapshot) bool

func (c *NamespacedOperatorCache) FindPreferred(comparer SnapshotComparer, p ...OperatorPredicate) []*Operator {
	snapshots := make([]*Snapshot, len(c.snapshots))
	for _, snapshot := range c.snapshots {
		snapshots = append(snapshots, snapshot.snapshot)
	}
	sort.Slice(snapshots, func(i, j int) bool {
		return false
	})
	var result []*Operator
	for _, snapshot := range snapshots {
		result = append(result, snapshot.Find(p...)...)
	}
	return result
}

// can pass in a Source at cache creation time instead of this
// func (c *NamespacedOperatorCache) WithExistingOperators(snapshot *CatalogSnapshot) MultiCatalogOperatorFinder {
// 	o := &NamespacedOperatorCache{
// 		namespaces: c.namespaces,
// 		existing:   &snapshot.key,
// 		snapshots:  c.snapshots,
// 	}
// 	o.snapshots[snapshot.key] = snapshot
// 	return o
// }

func (c *NamespacedOperatorCache) Find(p ...OperatorPredicate) []*Operator {
	return c.FindPreferred(nil, p...)
}

type Snapshot struct {
	Labels    map[string]string
	Operators []Operator
	Priority  int // label?
}

type snapshotHeader struct {
	key      string
	expiry   time.Time
	m        sync.RWMutex
	pop      context.CancelFunc
	err      error
	snapshot *Snapshot
}

func (s *snapshotHeader) Cancel() {
	s.pop()
}

func (s *snapshotHeader) Expired(at time.Time) bool {
	return !at.Before(s.expiry)
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

type SortableSnapshots []Snapshot

func NewSortableSnapshots(namespaces []string, snapshots map[registry.CatalogKey]*CatalogSnapshot) SortableSnapshots {
	sorted := SortableSnapshots{
		existing:   existing,
		preferred:  preferred,
		snapshots:  make([]*CatalogSnapshot, 0),
		namespaces: make(map[string]int, 0),
	}
	for i, n := range namespaces {
		sorted.namespaces[n] = i
	}
	for _, s := range snapshots {
		sorted.snapshots = append(sorted.snapshots, s)
	}
	return sorted
}

var _ sort.Interface = SortableSnapshots{}

// Len is the number of elements in the collection.
func (s SortableSnapshots) Len() int {
	return len(s.snapshots)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SortableSnapshots) Less(i, j int) bool {
	// existing operators are preferred over catalog operators
	if s.existing != nil &&
		s.snapshots[i].key.Name == s.existing.Name &&
		s.snapshots[i].key.Namespace == s.existing.Namespace {
		return true
	}
	if s.existing != nil &&
		s.snapshots[j].key.Name == s.existing.Name &&
		s.snapshots[j].key.Namespace == s.existing.Namespace {
		return false
	}

	// preferred catalog is less than all other catalogs
	if s.preferred != nil &&
		s.snapshots[i].key.Name == s.preferred.Name &&
		s.snapshots[i].key.Namespace == s.preferred.Namespace {
		return true
	}
	if s.preferred != nil &&
		s.snapshots[j].key.Name == s.preferred.Name &&
		s.snapshots[j].key.Namespace == s.preferred.Namespace {
		return false
	}

	// the rest are sorted first on priority, namespace and then by name
	if s.snapshots[i].priority != s.snapshots[j].priority {
		return s.snapshots[i].priority > s.snapshots[j].priority
	}
	if s.snapshots[i].key.Namespace != s.snapshots[j].key.Namespace {
		return s.namespaces[s.snapshots[i].key.Namespace] < s.namespaces[s.snapshots[j].key.Namespace]
	}

	return s.snapshots[i].key.Name < s.snapshots[j].key.Name
}

// Swap swaps the elements with indexes i and j.
func (s SortableSnapshots) Swap(i, j int) {
	s.snapshots[i], s.snapshots[j] = s.snapshots[j], s.snapshots[i]
}

func (s *Snapshot) Find(p ...OperatorPredicate) []*Operator {
	s.m.RLock()
	defer s.m.RUnlock()
	return Filter(s.Operators, p...)
}

type OperatorFinder interface {
	Find(...OperatorPredicate) []*Operator
}

type MultiCatalogOperatorFinder interface {
	Catalog(registry.CatalogKey) OperatorFinder
	FindPreferred(SnapshotComparer, ...OperatorPredicate) []*Operator
	//	WithExistingOperators(*CatalogSnapshot) MultiCatalogOperatorFinder
	Error() error
	OperatorFinder
}

type EmptyOperatorFinder struct{}

func (f EmptyOperatorFinder) Find(...OperatorPredicate) []*Operator {
	return nil
}

func AtLeast(n int, operators []*Operator) ([]*Operator, error) {
	if len(operators) < n {
		return nil, fmt.Errorf("expected at least %d operator(s), got %d", n, len(operators))
	}
	return operators, nil
}

func ExactlyOne(operators []*Operator) (*Operator, error) {
	if len(operators) != 1 {
		return nil, fmt.Errorf("expected exactly one operator, got %d", len(operators))
	}
	return operators[0], nil
}

func Filter(operators []*Operator, p ...OperatorPredicate) []*Operator {
	var result []*Operator
	for _, o := range operators {
		if Matches(o, p...) {
			result = append(result, o)
		}
	}
	return result
}

func Matches(o *Operator, p ...OperatorPredicate) bool {
	return And(p...).Test(o)
}
