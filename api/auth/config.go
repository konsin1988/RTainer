package auth

import (
	"fmt"
	"os"
)

type KeycloakConfig struct {
	KCBaseURL 	string
	KCRealm   	string
	KCClientId 	string
}

func LoadConfig() (*KeycloakConfig, error) {

	cfg := &KeycloakConfig{
		KCBaseURL: 	os.Getenv("KC_BASEURL"),
		KCRealm:   	os.Getenv("KC_REALM"),
		KCClientId: os.Getenv("KC_CLIENT_ID"),
	}

	if cfg.KCBaseURL == "" {
		return nil, fmt.Errorf("KC_BASE_URL is required")
	}

	if cfg.KCRealm == "" {
		return nil, fmt.Errorf("KC_REALM is required")
	}
	if cfg.KCClientId == "" {
		return nil, fmt.Errorf("KC_CLIENT_ID is required")
	}

	return cfg, nil
}
