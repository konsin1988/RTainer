package auth

import (
	"fmt"
	"log"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var jwks keyfunc.Keyfunc

func InitKeycloak() error {
	cfg,err := LoadConfig() 
	if err != nil {
		return err
	}
	jwksURL := fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/certs",
		cfg.KCBaseURL,
		cfg.KCRealm,
	) 
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return err
	}
	log.Println("Keycloak initialization done!!!")

	jwks = k
	return nil
}

func ValidateToken(tokenString string) (*User, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&KeycloakClaims{},
		jwks.Keyfunc,
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*KeycloakClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &User{
		ID:       claims.Subject,
		Username: claims.PreferredUsername,
		Email:    claims.Email,
		Roles:    claims.RealmAccess.Roles,
	}, nil
}
