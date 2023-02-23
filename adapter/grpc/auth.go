package grpc

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/forest33/warthog/business/entity"
)

type basicAuth struct {
	login    string
	password string
}

var symmetricAlgorithms = map[string]struct{}{
	"HS256": {},
	"HS384": {},
	"HS512": {},
}

func (c *Client) getAuth(auth *entity.Auth) (grpc.DialOption, error) {
	if auth == nil || auth.Type == entity.AuthTypeNone {
		return nil, nil
	}

	switch auth.Type {
	case entity.AuthTypeBasic:
		return c.authBasic(auth)
	case entity.AuthTypeBearer:
		return c.authBearer(auth)
	case entity.AuthTypeJWT:
		return c.authJWT(auth)
	case entity.AuthTypeGCE:
		return c.authGCE(auth)
	}

	return nil, nil
}

func (c *Client) authBasic(auth *entity.Auth) (grpc.DialOption, error) {
	return grpc.WithPerRPCCredentials(basicAuth{
		login:    auth.Login,
		password: auth.Password,
	}), nil
}

func (c *Client) authBearer(auth *entity.Auth) (grpc.DialOption, error) {
	token := &oauth2.Token{
		AccessToken: auth.Token,
	}
	token.TokenType = auth.HeaderPrefix

	return grpc.WithPerRPCCredentials(
		oauth.NewOauthAccess(token),
	), nil
}

func (c *Client) authJWT(auth *entity.Auth) (grpc.DialOption, error) {
	signingMethod := jwt.GetSigningMethod(auth.Algorithm)
	if signingMethod == nil {
		return nil, fmt.Errorf("unknown signing algorithm: %s", auth.Algorithm)
	}

	var (
		secret interface{}
		err    error
	)

	if _, ok := symmetricAlgorithms[auth.Algorithm]; ok {
		secret = []byte(auth.Secret)
		if auth.SecretBase64 {
			secret, err = base64.StdEncoding.DecodeString(auth.Secret)
			if err != nil {
				return nil, err
			}
		}
	} else {
		secret, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(auth.PrivateKey))
		if err != nil {
			return nil, err
		}
	}

	jwtToken := jwt.NewWithClaims(signingMethod, jwt.MapClaims(auth.Payload))
	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	return c.authBearer(&entity.Auth{
		Token:        token,
		HeaderPrefix: auth.HeaderPrefix,
	})
}

func (c *Client) authGCE(auth *entity.Auth) (grpc.DialOption, error) {
	perRPC, err := oauth.NewServiceAccountFromKey([]byte(auth.GoogleToken), auth.GoogleScopes...)
	if err != nil {
		return nil, err
	}
	return grpc.WithPerRPCCredentials(perRPC), nil
}

func (b basicAuth) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	auth := b.login + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}
