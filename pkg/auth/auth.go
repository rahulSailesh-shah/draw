package auth

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)


func LoadKeys(jwksURL string) (jwk.Set, error) {
	ctx:= context.Background()
	ks, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, err
	}
	return ks, nil
}

func UserFromToken(r *http.Request, authKeys jwk.Set) (string, error) {
	token, err := jwt.ParseRequest(r, jwt.WithKeySet(authKeys))
	if err != nil {
		return "", err
	}
	userId, exists := token.Subject()
	if !exists {
		return "", nil
	}
	return userId, nil
}
