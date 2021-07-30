// +todo
package cache

import (
	"context"
	"encoding/json"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/api"
	"github.com/operator-framework/operator-registry/pkg/client"
	opregistry "github.com/operator-framework/operator-registry/pkg/registry"
)

func SourceFromCatalogSource(cs *v1alpha1.CatalogSource, c client.Client) Source {
	return nil //todo
}

func populate(ctx context.Context, hdr *snapshotHeader, registry client.Client) (*Snapshot, error) {
	var result Snapshot

	// Fetching default channels this way makes many round trips
	// -- may need to either add a new API to fetch all at once,
	// or embed the information into Bundle.
	defaultChannels := make(map[string]string)

	it, err := registry.ListBundles(ctx)
	if err != nil {
		//snapshot.logger.Errorf("failed to list bundles: %s", err.Error())
		return nil, err
	}
	//c.logger.WithField("catalog", snapshot.key.String()).Debug("updating cache")
	for b := it.Next(); b != nil; b = it.Next() {
		defaultChannel, ok := defaultChannels[b.PackageName]
		if !ok {
			if p, err := registry.GetPackage(ctx, b.PackageName); err != nil {
				//snapshot.logger.Warnf("failed to retrieve default channel for bundle, continuing: %v", err)
				continue
			} else {
				defaultChannels[b.PackageName] = p.DefaultChannelName
				defaultChannel = p.DefaultChannelName
			}
		}
		o, err := NewOperatorFromBundle(b, "", 1, defaultChannel) // todo
		if err != nil {
			//snapshot.logger.Warnf("failed to construct operator from bundle, continuing: %v", err)
			continue
		}
		o.providedAPIs = o.ProvidedAPIs().StripPlural()
		o.requiredAPIs = o.RequiredAPIs().StripPlural()
		o.replaces = b.Replaces
		ensurePackageProperty(o, b.PackageName, b.Version)
		result.Operators = append(result.Operators, *o)
	}
	if err := it.Error(); err != nil {
		//snapshot.logger.Warnf("error encountered while listing bundles: %s", err.Error())
		return nil, err
	}
	snapshot.operators = operators
}

func ensurePackageProperty(o *Operator, name, version string) {
	for _, p := range o.Properties() {
		if p.Type == opregistry.PackageType {
			return
		}
	}
	prop := opregistry.PackageProperty{
		PackageName: name,
		Version:     version,
	}
	bytes, err := json.Marshal(prop)
	if err != nil {
		return
	}
	o.properties = append(o.properties, &api.Property{
		Type:  opregistry.PackageType,
		Value: string(bytes),
	})
}
