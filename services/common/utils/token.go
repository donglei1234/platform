package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type tokenPayload struct {
	Exp int    `json:"exp"`
	UID string `json:"uid"`
}

const JwtToken = "token"

func DecodeToken(token string) (string, error) {
	if token == "" {
		return "", ErrNull
	}
	splitToken := strings.Split(token, ".")
	if len(splitToken) <= 1 {
		return "", ErrTokenInvalid
	}

	decodeString, err := base64.RawURLEncoding.DecodeString(splitToken[1])
	if err != nil {
		return "", ErrDecodeToken
	}
	var payload tokenPayload
	err = json.Unmarshal(decodeString, &payload)
	if err != nil {
		return "", ErrUnmarshalToken
	}
	return payload.UID, nil

}

// 从上下文MeteData中解析jwtToken
func ParseJwtToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetadataDataLost
	}

	var jwtToken []string
	if jwtToken, ok = md[JwtToken]; !ok {
		return "", errors.Wrap(ErrMetadataDataLost, JwtToken)
	}

	return strings.Join(jwtToken, ""), nil
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("token", token)
	return metadata.NewOutgoingContext(context.Background(), md)
}
