package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	basicLogin    = "login"
	basicPassword = "password"
	bearerToken   = "some-secret-token"
	jwtSecret     = "secret"

	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

func ensureValidBasicCredentials(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMissingMetadata
	}
	if !validBasic(md["authorization"]) {
		return errInvalidToken
	}
	return nil
}

func validBasic(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")
	return token == base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", basicLogin, basicPassword)))
}

func validToken(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == bearerToken
}

func validJWT(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	rawToken := strings.TrimPrefix(authorization[0], "Bearer ")

	_, err := jwt.ParseWithClaims(rawToken, &jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

	return err == nil
}

func ensureValidToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMissingMetadata
	}
	if !validToken(md["authorization"]) {
		return errInvalidToken
	}
	return nil
}

func ensureValidJWT(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMissingMetadata
	}
	if !validJWT(md["authorization"]) {
		return errInvalidToken
	}
	return nil
}
