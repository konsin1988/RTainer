package auth

import "github.com/golang-jwt/jwt/v5"

type KeycloakClaims struct {
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`

	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`

	jwt.RegisteredClaims
}

