package config_store

import (
	"github.com/tflz/aws-landing-zone/config_store/organization_config"
)

type ConfigStore struct {
	OrganizationConfig *organization_config.OrganizationConfig
}

func NewConfigStore() (*ConfigStore, error) {
	org_config, err := organization_config.LoadOrganizationConfig()
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		OrganizationConfig: org_config,
	}, nil
}
