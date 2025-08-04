package organization_config

import (
	"github.com/zclconf/go-cty/cty"
)

type OrganizationConfig struct {
	HomeRegion          string
	OrganizationalUnits []OrganizationalUnit
}

type OrganizationalUnit struct {
	Alias  string
	Name   string
	Parent OrganizationalUnitReference
}

type OrganizationalUnitReference interface {
	IsOrganizationalUnitReference() bool
}

type RootOrganizationalUnitReference struct{}

func (ref RootOrganizationalUnitReference) IsOrganizationalUnitReference() bool {
	return true
}

type AliasOrganizationalUnitReference struct {
	Alias string
}

func (ref AliasOrganizationalUnitReference) IsOrganizationalUnitReference() bool {
	return true
}

type RawOrganizationalUnitReference struct {
	ParentId string
}

func (ref RawOrganizationalUnitReference) IsOrganizationalUnitReference() bool {
	return true
}

func LoadOrganizationConfig() (*OrganizationConfig, error) {
	rawConfig, err := parse("organization.tflz")
	if err != nil {
		return nil, err
	}

	rawConfig.ctx = rawConfig.generateContext()

	organizational_units, err := rawConfig.getOrganizationalUnits()
	if err != nil {
		return nil, err
	}

	return &OrganizationConfig{
		HomeRegion:          rawConfig.HomeRegion,
		OrganizationalUnits: organizational_units,
	}, nil
}

func (rawConfig *rawOrganizationConfig) getOrganizationalUnits() ([]OrganizationalUnit, error) {
	// Convert raw configuration to OrganizationConfig
	units := make([]OrganizationalUnit, len(rawConfig.OrganizationalUnits))
	for i, rawOU := range rawConfig.OrganizationalUnits {
		rawParent, err := rawOU.Parent.Value(&rawConfig.ctx)
		if err != nil {
			return nil, err
		}

		var parentRef OrganizationalUnitReference
		if rawParent.Type() == cty.String {
			parentRef = RawOrganizationalUnitReference{
				ParentId: rawParent.AsString(),
			}
		} else if rawParent.Type().IsObjectType() && rawParent.Type().HasAttribute("__ref") {
			if rawParent.GetAttr("__ref").AsString() == "organizational_unit" {
				parentRef = AliasOrganizationalUnitReference{
					Alias: rawParent.GetAttr("alias").AsString(),
				}
			} else if rawParent.GetAttr("__ref").AsString() == "root" {
				parentRef = RootOrganizationalUnitReference{}
			}
		}

		units[i] = OrganizationalUnit{
			Alias:  rawOU.Alias,
			Name:   rawOU.Name,
			Parent: parentRef,
		}
	}

	return units, nil
}
