package resolver

import (
	"encoding/json"

	"github.com/benluddy/resolver/cache"
	"github.com/operator-framework/operator-registry/pkg/api"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func apiSetToDependencies(crds, apis cache.APISet) (out []*api.Dependency) {
	if len(crds)+len(apis) == 0 {
		return nil
	}
	out = make([]*api.Dependency, 0)
	for a := range crds {
		val, err := json.Marshal(registry.GVKDependency{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Dependency{
			Type:  registry.GVKType,
			Value: string(val),
		})
	}
	for a := range apis {
		val, err := json.Marshal(registry.GVKDependency{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Dependency{
			Type:  registry.GVKType,
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
		val, err := json.Marshal(registry.GVKProperty{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  registry.GVKType,
			Value: string(val),
		})
	}
	for a := range apis {
		val, err := json.Marshal(registry.GVKProperty{
			Group:   a.Group,
			Kind:    a.Kind,
			Version: a.Version,
		})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  registry.GVKType,
			Value: string(val),
		})
	}
	if deprecated {
		val, err := json.Marshal(registry.DeprecatedProperty{})
		if err != nil {
			panic(err)
		}
		out = append(out, &api.Property{
			Type:  registry.DeprecatedType,
			Value: string(val),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return
}
