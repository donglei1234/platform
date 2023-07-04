package db

import (
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
)

const (
	ProfileBanTheme = "profileban"
	BulletinTheme   = "bulletin"
)

func makeProfileBanPath(path string) (keys.Key, error) {
	return keys.NewKeyFromParts(ProfileBanTheme, path)
}

func makeProfileBanKey(profileId string) (keys.Key, error) {
	return keys.NewKeyFromParts(ProfileBanTheme, profileId)
}

func makeBulletinPath() (keys.Key, error) {
	return keys.NewKeyFromParts(BulletinTheme, "infos")
}

func makeBulletinKey(id string) (keys.Key, error) {
	return keys.NewKeyFromParts(BulletinTheme, id)
}
