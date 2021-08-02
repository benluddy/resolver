package cache

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/blang/semver/v4"

	"github.com/operator-framework/operator-registry/pkg/api"
	opregistry "github.com/operator-framework/operator-registry/pkg/registry"
)

type Predicate interface {
	Test(*Operator) bool
	String() string
}

type csvNamePredicate string

func CSVNamePredicate(name string) Predicate {
	return csvNamePredicate(name)
}

func (c csvNamePredicate) Test(o *Operator) bool {
	return o.Name == string(c)
}

func (c csvNamePredicate) String() string {
	return fmt.Sprintf("with name: %s", string(c))
}

type channelPredicate string

func ChannelPredicate(channel string) Predicate {
	return channelPredicate(channel)
}

func (ch channelPredicate) Test(o *Operator) bool {
	// all operators match the empty channel
	return ch == "" || o.Channel == string(ch)
}

func (ch channelPredicate) String() string {
	return fmt.Sprintf("with channel: %s", string(ch))
}

type pkgPredicate string

func PkgPredicate(pkg string) Predicate {
	return pkgPredicate(pkg)
}

func (pkg pkgPredicate) Test(o *Operator) bool {
	for _, p := range o.Properties {
		if p.Type != opregistry.PackageType {
			continue
		}
		var prop opregistry.PackageProperty
		err := json.Unmarshal([]byte(p.Value), &prop)
		if err != nil {
			continue
		}
		if prop.PackageName == string(pkg) {
			return true
		}
	}
	return o.Package == string(pkg)
}

func (pkg pkgPredicate) String() string {
	return fmt.Sprintf("with package: %s", string(pkg))
}

type versionInRangePredicate struct {
	r   semver.Range
	str string
}

func VersionInRangePredicate(r semver.Range, version string) Predicate {
	return versionInRangePredicate{r: r, str: version}
}

func (v versionInRangePredicate) Test(o *Operator) bool {
	for _, p := range o.Properties {
		if p.Type != opregistry.PackageType {
			continue
		}
		var prop opregistry.PackageProperty
		err := json.Unmarshal([]byte(p.Value), &prop)
		if err != nil {
			continue
		}
		ver, err := semver.Parse(prop.Version)
		if err != nil {
			continue
		}
		if v.r(ver) {
			return true
		}
	}
	return o.Version != nil && v.r(*o.Version)
}

func (v versionInRangePredicate) String() string {
	return fmt.Sprintf("with version in range: %v", v.str)
}

type labelPredicate string

func LabelPredicate(label string) Predicate {
	return labelPredicate(label)
}
func (l labelPredicate) Test(o *Operator) bool {
	for _, p := range o.Properties {
		if p.Type != opregistry.LabelType {
			continue
		}
		var prop opregistry.LabelProperty
		err := json.Unmarshal([]byte(p.Value), &prop)
		if err != nil {
			continue
		}
		if prop.Label == string(l) {
			return true
		}
	}
	return false
}

func (l labelPredicate) String() string {
	return fmt.Sprintf("with label: %v", string(l))
}

type sourcePredicate SourceKey

func SourcePredicate(key SourceKey) Predicate {
	return sourcePredicate(key)
}

func (c sourcePredicate) Test(o *Operator) bool {
	return SourceKey(c) == o.SourceKey
}

func (c sourcePredicate) String() string {
	// todo: key is opaque, but this needs to translate back to a catsrc ns/name
	return fmt.Sprintf("from source with key : %v", SourceKey(c))
}

type gvkPredicate struct {
	api opregistry.APIKey
}

func ProvidingAPIPredicate(api opregistry.APIKey) Predicate {
	return gvkPredicate{
		api: api,
	}
}

func (g gvkPredicate) Test(o *Operator) bool {
	for _, p := range o.Properties {
		if p.Type != opregistry.GVKType {
			continue
		}
		var prop opregistry.GVKProperty
		err := json.Unmarshal([]byte(p.Value), &prop)
		if err != nil {
			continue
		}
		if prop.Kind == g.api.Kind && prop.Version == g.api.Version && prop.Group == g.api.Group {
			return true
		}
	}
	return false
}

func (g gvkPredicate) String() string {
	return fmt.Sprintf("providing an API with group: %s, version: %s, kind: %s", g.api.Group, g.api.Version, g.api.Kind)
}

type skipRangeIncludesPredication struct {
	version semver.Version
}

func SkipRangeIncludesPredicate(version semver.Version) Predicate {
	return skipRangeIncludesPredication{version: version}
}

func (s skipRangeIncludesPredication) Test(o *Operator) bool {
	return o.SkipRange(s.version)
}

func (s skipRangeIncludesPredication) String() string {
	return fmt.Sprintf("skip range includes: %v", s.version.String())
}

type replacesPredicate string

func ReplacesPredicate(replaces string) Predicate {
	return replacesPredicate(replaces)
}

func (r replacesPredicate) Test(o *Operator) bool {
	if o.Replaces == string(r) {
		return true
	}
	for _, s := range o.Skips {
		if s == string(r) {
			return true
		}
	}
	return false
}

func (r replacesPredicate) String() string {
	return fmt.Sprintf("replaces: %v", string(r))
}

type andPredicate struct {
	predicates []Predicate
}

func And(p ...Predicate) Predicate {
	return andPredicate{
		predicates: p,
	}
}

func (p andPredicate) Test(o *Operator) bool {
	for _, predicate := range p.predicates {
		if predicate.Test(o) == false {
			return false
		}
	}
	return true
}

func (p andPredicate) String() string {
	var b bytes.Buffer
	for i, predicate := range p.predicates {
		b.WriteString(predicate.String())
		if i != len(p.predicates)-1 {
			b.WriteString(" and ")
		}
	}
	return b.String()
}

func Or(p ...Predicate) Predicate {
	return orPredicate{
		predicates: p,
	}
}

type orPredicate struct {
	predicates []Predicate
}

func (p orPredicate) Test(o *Operator) bool {
	for _, predicate := range p.predicates {
		if predicate.Test(o) == true {
			return true
		}
	}
	return false
}

func (p orPredicate) String() string {
	var b bytes.Buffer
	for i, predicate := range p.predicates {
		b.WriteString(predicate.String())
		if i != len(p.predicates)-1 {
			b.WriteString(" or ")
		}
	}
	return b.String()
}

type booleanPredicate struct {
	result bool
}

func BooleanPredicate(result bool) Predicate {
	return booleanPredicate{result: result}
}

func (b booleanPredicate) Test(o *Operator) bool {
	return b.result
}

func (b booleanPredicate) String() string {
	if b.result {
		return fmt.Sprintf("predicate is true")
	}
	return fmt.Sprintf("predicate is false")
}

func True() Predicate {
	return BooleanPredicate(true)
}

func False() Predicate {
	return BooleanPredicate(false)
}

type countingPredicate struct {
	p Predicate
	n *int
}

func (c countingPredicate) Test(o *Operator) bool {
	if c.p.Test(o) {
		*c.n++
		return true
	}
	return false
}
func (c countingPredicate) String() string {
	return c.p.String()
}
func CountingPredicate(p Predicate, n *int) Predicate {
	return countingPredicate{p: p, n: n}
}

func PredicateForProperty(property *api.Property) (Predicate, error) {
	if property == nil {
		return nil, nil
	}
	p, ok := predicates[property.Type]
	if !ok {
		return nil, nil
	}
	return p(property.Value)
}

var predicates = map[string]func(string) (Predicate, error){
	"olm.gvk.required":     predicateForRequiredGVKProperty,
	"olm.package.required": predicateForRequiredPackageProperty,
	"olm.label.required":   predicateForRequiredLabelProperty,
}

func predicateForRequiredGVKProperty(value string) (Predicate, error) {
	var gvk struct {
		Group   string `json:"group"`
		Version string `json:"version"`
		Kind    string `json:"kind"`
	}
	if err := json.Unmarshal([]byte(value), &gvk); err != nil {
		return nil, err
	}
	return ProvidingAPIPredicate(opregistry.APIKey{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}), nil
}

func predicateForRequiredPackageProperty(value string) (Predicate, error) {
	var pkg struct {
		PackageName  string `json:"packageName"`
		VersionRange string `json:"versionRange"`
	}
	if err := json.Unmarshal([]byte(value), &pkg); err != nil {
		return nil, err
	}
	ver, err := semver.ParseRange(pkg.VersionRange)
	if err != nil {
		return nil, err
	}
	return And(PkgPredicate(pkg.PackageName), VersionInRangePredicate(ver, pkg.VersionRange)), nil
}

func predicateForRequiredLabelProperty(value string) (Predicate, error) {
	var label struct {
		Label string `json:"label"`
	}
	if err := json.Unmarshal([]byte(value), &label); err != nil {
		return nil, err
	}
	return LabelPredicate(label.Label), nil
}
