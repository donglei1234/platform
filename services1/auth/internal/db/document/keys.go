package db

import (
	"github.com/donglei1234/platform/services/common/nosql/document"
)

func userAuthenticationKey(appId, userId string) (document.Key, error) {
	return document.NewKeyFromParts(appId, "users", userId)
}

func userPlatformAuthenticationKey(appId, platformType, platformUid string) (document.Key, error) {
	return document.NewKeyFromParts(appId, platformType, "users", platformUid)
}
