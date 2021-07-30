package resolver

import (
	"encoding/json"

	"github.com/operator-framework/operator-registry/pkg/api"
	"github.com/benluddy/resolver/cache"
)

func apiSetToDependencies(crds, apis cache.APISet) (out []*api.Dependency) {
	if len(crds)+len(apis) == 0 {
		return nil
	}
	out = make([]*api.Dependency, 0)
	for a := range crds {
		val, err := json.Marshal(opregistry.GVKDependency{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Dependency{
			Type:  opregistry.GVKType,
			Value: string(val),
		})
	}
	for a := range apis {
		val, err := json.Marshal(opregistry.GVKDependency{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Dependency{
			Type:  opregistry.GVKType,
			Value: string(val),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return
}

func apiSetToProperties(crds, apis cache.APISet, deprecated bool) (out []*api.Property) {
	out = make([]*api.Property, 0)
	for a := range crds {
		val, err := json.Marshal(opregistry.GVKProperty{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  opregistry.GVKType,
			Value: string(val),
		})
	}
	for a := range apis {
		val, err := json.Marshal(opregistry.GVKProperty{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  opregistry.GVKType,
			Value: string(val),
		})
	}
	if deprecated {
		val, err := json.Marshal(opregistry.DeprecatedProperty{})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  opregistry.DeprecatedType,
			Value: string(val),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return
}

