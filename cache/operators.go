package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/api"
	opregistry "github.com/operator-framework/operator-registry/pkg/registry"
)

type APISet map[opregistry.APIKey]struct{}

func EmptyAPISet() APISet {
	return map[opregistry.APIKey]struct{}{}
}

func (s APISet) PopAPIKey() *opregistry.APIKey {
	for a := range s {
		api := &opregistry.APIKey{
			Group:   a.Group,
			Version: a.Version,
			Kind:    a.Kind,
			Plural:  a.Plural,
		}
		delete(s, a)
		return api
	}
	return nil
}

func GVKStringToProvidedAPISet(gvksStr string) APISet {
	set := make(APISet)
	// TODO: Should we make gvk strings lowercase to avoid issues with user set gvks?
	gvks := strings.Split(strings.Replace(gvksStr, " ", "", -1), ",")
	for _, gvkStr := range gvks {
		gvk, _ := schema.ParseKindArg(gvkStr)
		if gvk != nil {
			set[opregistry.APIKey{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind}] = struct{}{}
		}
	}

	return set
}

func APIKeyToGVKString(key opregistry.APIKey) string {
	// TODO: Add better validation of GVK
	return strings.Join([]string{key.Kind, key.Version, key.Group}, ".")
}

func APIKeyToGVKHash(key opregistry.APIKey) (string, error) {
	hash := fnv.New64a()
	if _, err := hash.Write([]byte(APIKeyToGVKString(key))); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum64()), nil
}

func (s APISet) String() string {
	gvkStrs := make([]string, len(s))
	i := 0
	for api := range s {
		// TODO: Only add valid GVK strings
		gvkStrs[i] = APIKeyToGVKString(api)
		i++
	}
	sort.Strings(gvkStrs)

	return strings.Join(gvkStrs, ",")
}

// Union returns the union of the APISet and the given list of APISets
func (s APISet) Union(sets ...APISet) APISet {
	union := make(APISet)
	for api := range s {
		union[api] = struct{}{}
	}
	for _, set := range sets {
		for api := range set {
			union[api] = struct{}{}
		}
	}

	return union
}

// Intersection returns the intersection of the APISet and the given list of APISets
func (s APISet) Intersection(sets ...APISet) APISet {
	intersection := make(APISet)
	for _, set := range sets {
		for api := range set {
			if _, ok := s[api]; ok {
				intersection[api] = struct{}{}
			}
		}
	}

	return intersection
}

func (s APISet) Difference(set APISet) APISet {
	difference := make(APISet).Union(s)
	for api := range set {
		if _, ok := difference[api]; ok {
			delete(difference, api)
		}
	}

	return difference
}

// IsSubset returns true if the APISet is a subset of the given one
func (s APISet) IsSubset(set APISet) bool {
	for api := range s {
		if _, ok := set[api]; !ok {
			return false
		}
	}

	return true
}

// StripPlural returns the APISet with the Plural field of all APIKeys removed
func (s APISet) StripPlural() APISet {
	set := make(APISet)
	for api := range s {
		set[opregistry.APIKey{Group: api.Group, Version: api.Version, Kind: api.Kind}] = struct{}{}
	}

	return set
}

type APIOwnerSet map[opregistry.APIKey]*Operator

func EmptyAPIOwnerSet() APIOwnerSet {
	return map[opregistry.APIKey]*Operator{}
}

type OperatorSet map[string]*Operator

func EmptyOperatorSet() OperatorSet {
	return map[string]*Operator{}
}

// Snapshot returns a new set, pointing to the same values
func (o OperatorSet) Snapshot() OperatorSet {
	out := make(map[string]*Operator)
	for key, val := range o {
		out[key] = val
	}
	return out
}

type APIMultiOwnerSet map[opregistry.APIKey]OperatorSet

func EmptyAPIMultiOwnerSet() APIMultiOwnerSet {
	return map[opregistry.APIKey]OperatorSet{}
}

func (s APIMultiOwnerSet) PopAPIKey() *opregistry.APIKey {
	for a := range s {
		api := &opregistry.APIKey{
			Group:   a.Group,
			Version: a.Version,
			Kind:    a.Kind,
			Plural:  a.Plural,
		}
		delete(s, a)
		return api
	}
	return nil
}

func (s APIMultiOwnerSet) PopAPIRequirers() OperatorSet {
	requirers := EmptyOperatorSet()
	for a := range s {
		for key, op := range s[a] {
			requirers[key] = op
		}
		delete(s, a)
		return requirers
	}
	return nil
}

// sourceinfo -> SourceLabels
// type Origin struct {
// 	Package        string
// 	Channel        string
// 	StartingCSV    string
// 	SourceKey      SourceKey
// 	DefaultChannel bool
// 	Subscription   *v1alpha1.Subscription
// }

// func (i *Origin) String() string {
// 	// todo: is this user-facing? source key is opaque
// 	return fmt.Sprintf("%s/%s in %v", i.Package, i.Channel, i.SourceKey)
// }

type Operator struct {
	Name       string
	Properties []*api.Property

	// todo: read these from properties, synthesize properties for now
	Replaces  string
	Skips     []string
	SkipRange semver.Range
	Channel   string
	Package   string
	Version   *semver.Version

	SourceKey      SourceKey
	SourceLabels   map[string]string
	DefaultChannel bool
	Subscription   *v1alpha1.Subscription
	Bundle         *api.Bundle // TODO
}

func NewOperatorFromV1Alpha1CSV(csv *v1alpha1.ClusterServiceVersion) (*Operator, error) {
	providedAPIs := EmptyAPISet()
	for _, crdDef := range csv.Spec.CustomResourceDefinitions.Owned {
		parts := strings.SplitN(crdDef.Name, ".", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("error parsing crd name: %s", crdDef.Name)
		}
		providedAPIs[opregistry.APIKey{Plural: parts[0], Group: parts[1], Version: crdDef.Version, Kind: crdDef.Kind}] = struct{}{}
	}
	for _, api := range csv.Spec.APIServiceDefinitions.Owned {
		providedAPIs[opregistry.APIKey{Group: api.Group, Version: api.Version, Kind: api.Kind, Plural: api.Name}] = struct{}{}
	}

	requiredAPIs := EmptyAPISet()
	for _, crdDef := range csv.Spec.CustomResourceDefinitions.Required {
		parts := strings.SplitN(crdDef.Name, ".", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("error parsing crd name: %s", crdDef.Name)
		}
		requiredAPIs[opregistry.APIKey{Plural: parts[0], Group: parts[1], Version: crdDef.Version, Kind: crdDef.Kind}] = struct{}{}
	}
	for _, api := range csv.Spec.APIServiceDefinitions.Required {
		requiredAPIs[opregistry.APIKey{Group: api.Group, Version: api.Version, Kind: api.Kind, Plural: api.Name}] = struct{}{}
	}

	properties, err := providedAPIsToProperties(providedAPIs)
	if err != nil {
		return nil, err
	}
	dependencies, err := requiredAPIsToProperties(requiredAPIs)
	if err != nil {
		return nil, err
	}
	properties = append(properties, dependencies...)

	return &Operator{
		Name:       csv.GetName(),
		Version:    &csv.Spec.Version.Version,
		Properties: properties,
	}, nil
}

func (o *Operator) Inline() bool {
	return o.Bundle != nil && o.Bundle.GetBundlePath() == ""
}

func (o *Operator) DependencyPredicates() (predicates []Predicate, err error) {
	predicates = make([]Predicate, 0)
	for _, property := range o.Properties {
		predicate, err := PredicateForProperty(property)
		if err != nil {
			return nil, err
		}
		if predicate == nil {
			continue
		}
		predicates = append(predicates, predicate)
	}
	return
}

func requiredAPIsToProperties(apis APISet) (out []*api.Property, err error) {
	if len(apis) == 0 {
		return
	}
	out = make([]*api.Property, 0)
	for a := range apis {
		val, err := json.Marshal(struct {
			Group   string `json:"group"`
			Version string `json:"version"`
			Kind    string `json:"kind"`
		}{
			Group:   a.Group,
			Version: a.Version,
			Kind:    a.Kind,
		})
		if err != nil {
			return nil, err
		}
		out = append(out, &api.Property{
			Type:  "olm.gvk.required",
			Value: string(val),
		})
	}
	if len(out) > 0 {
		return
	}
	return nil, nil
}

func providedAPIsToProperties(apis APISet) (out []*api.Property, err error) {
	out = make([]*api.Property, 0)
	for a := range apis {
		val, err := json.Marshal(opregistry.GVKProperty{
			Group:   a.Group,
			Version: a.Version,
			Kind:    a.Kind,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  opregistry.GVKType,
			Value: string(val),
		})
	}
	if len(out) > 0 {
		return
	}
	return nil, nil
}

func legacyDependenciesToProperties(dependencies []*api.Dependency) ([]*api.Property, error) {
	var result []*api.Property
	for _, dependency := range dependencies {
		switch dependency.Type {
		case "olm.gvk":
			type gvk struct {
				Group   string `json:"group"`
				Version string `json:"version"`
				Kind    string `json:"kind"`
			}
			var vfrom gvk
			if err := json.Unmarshal([]byte(dependency.Value), &vfrom); err != nil {
				return nil, fmt.Errorf("failed to unmarshal legacy 'olm.gvk' dependency: %w", err)
			}
			vto := gvk{
				Group:   vfrom.Group,
				Version: vfrom.Version,
				Kind:    vfrom.Kind,
			}
			vb, err := json.Marshal(&vto)
			if err != nil {
				return nil, fmt.Errorf("unexpected error marshaling generated 'olm.package.required' property: %w", err)
			}
			result = append(result, &api.Property{
				Type:  "olm.gvk.required",
				Value: string(vb),
			})
		case "olm.package":
			var vfrom struct {
				PackageName  string `json:"packageName"`
				VersionRange string `json:"version"`
			}
			if err := json.Unmarshal([]byte(dependency.Value), &vfrom); err != nil {
				return nil, fmt.Errorf("failed to unmarshal legacy 'olm.package' dependency: %w", err)
			}
			vto := struct {
				PackageName  string `json:"packageName"`
				VersionRange string `json:"versionRange"`
			}{
				PackageName:  vfrom.PackageName,
				VersionRange: vfrom.VersionRange,
			}
			vb, err := json.Marshal(&vto)
			if err != nil {
				return nil, fmt.Errorf("unexpected error marshaling generated 'olm.package.required' property: %w", err)
			}
			result = append(result, &api.Property{
				Type:  "olm.package.required",
				Value: string(vb),
			})
		case "olm.label":
			result = append(result, &api.Property{
				Type:  "olm.label.required",
				Value: dependency.Value,
			})
		}
	}
	return result, nil
}
