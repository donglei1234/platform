package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	pb "github.com/donglei1234/platform/services/proto/gen/auth/api"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/oauth2/v1"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	fb "github.com/huandu/facebook"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
)

type FacebookUserInfo struct {
	Data struct {
		AppID               string `json:"app_id"`
		Type                string `json:"type"`
		Application         string `json:"application"`
		DataAccessExpiresAt int    `json:"data_access_expires_at"`
		ExpiresAt           int    `json:"expires_at"`
		IsValid             bool   `json:"is_valid"`
		Metadata            struct {
			AuthType string `json:"auth_type"`
		} `json:"metadata"`
		Scopes []string `json:"scopes"`
		UserID string   `json:"user_id"`
	} `json:"data"`
}

const DebugToken = "/debug_token"
const InputToken = "input_token"
const ProfileId = "profileid"
const JwtToken = "token"

// 生成一个uuid
func GenerateUUID() (string, error) {
	random, err := uuid.NewRandom()
	id := strconv.Itoa((int)(random.ID()))
	if err != nil {
		return "", err
	}
	return id, nil
}

// 从上下文MeteData中解析ProfileId
func ParseProfileId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetadataDataLost
	}

	var profileId []string
	if profileId, ok = md[ProfileId]; !ok {
		return "", errors.Wrap(ErrMetadataDataLost, ProfileId)
	}

	return strings.Join(profileId, ""), nil
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

// deprecated
func GetFacebookId(appId, secret, accessToken string) (string, error) {
	app := fb.New(appId, secret)
	token := app.AppAccessToken()

	session := app.Session(token)
	resp, err := session.Get(DebugToken, fb.Params{
		InputToken: accessToken,
	})
	if err != nil {
		return "", err
	}

	var userInfo FacebookUserInfo
	err = resp.Decode(&userInfo)
	if err != nil {
		return "", err
	}
	if !userInfo.Data.IsValid {
		return "", err
	}

	return userInfo.Data.UserID, nil
}

// deprecated
func GetGoogleId(apiKeys, idToken string) (string, error) {
	oauth2Service, err := oauth2.NewService(nil, option.WithAPIKey(apiKeys))
	if err != nil {
		return "", err
	}
	tokenInfo := oauth2Service.Tokeninfo()
	token := tokenInfo.IdToken(idToken)
	do, err := token.Do()
	if err != nil {
		return "", err
	}
	// TODO 需要对不同的app验证Audience字段，如tata为743906483241-idjfg26hdialusrrf10v7l3ddkpf32jk.apps.googleusercontent.com
	// 参考：https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
	if do.Audience == "" {

	}
	if do.ExpiresIn <= 0 {
		return "", err
	}
	return do.UserId, nil
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

// CreateRoomToken 生成一个JwtToken，包含roomId和ids
func CreateRoomToken(roomInfo *pb.RoomInfo, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roomId":   roomInfo.RoomId,
		"roomType": roomInfo.RoomType,
		"ids":      roomInfo.Members,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", ErrSignedString
	}
	return token, nil
}

// CreatProfileJwt 生成一个JwtToken，包含uid
func CreatProfileJwt(profileId string, key string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": profileId,
		"exp": time.Now().Add(6 * time.Hour).Unix(),
	})
	token, err := at.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseRoomToken 从Jwt中解析Token
func ParseRoomToken(token string, key string) (*pb.RoomInfo, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningMethod
		}
		return []byte(key), nil
	})

	if claim.Valid {
		if claims, ok := claim.Claims.(jwt.MapClaims); ok {
			mems := claims["ids"].([]interface{})
			members := make([]string, len(mems))
			for k, v := range mems {
				members[k] = v.(string)
			}
			roomTypeInfo := claims["roomType"].(float64)
			roomType := int32(roomTypeInfo)

			return &pb.RoomInfo{
				RoomId:   claims["roomId"].(string),
				Members:  members,
				RoomType: roomType,
			}, nil
		}
		return nil, ErrTokenHandle
	}
	return nil, err
}

// ParseProfileToken 从Jwt中解析Token
func ParseProfileToken(token string, key string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningMethod
		}
		return []byte(key), nil
	})

	if claim.Valid {
		if claims, ok := claim.Claims.(jwt.MapClaims); ok {
			return claims["uid"].(string), nil
		}
		return "", ErrTokenHandle
	}
	return "", err
}

type tokenPayload struct {
	Exp int    `json:"exp"`
	UID string `json:"uid"`
}

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
