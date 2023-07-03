package db

import (
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"strconv"
)

const (
	MailTheme       = "mail"
	MailPublicIndex = "index"
)

func makeMailKey(profileId string) (keys.Key, error) {
	return keys.NewKeyFromParts(MailTheme, profileId)
}

func makeMailPublicIndexKey(profileId string) (keys.Key, error) {
	return keys.NewKeyFromParts(MailTheme, MailPublicIndex, profileId)
}

func makeFieldMailKey(uid int64) (keys.Key, error) {
	idStr := strconv.FormatInt(uid, 10)
	return keys.NewKeyFromParts(MailTheme, idStr)
}
