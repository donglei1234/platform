package db

import (
	"github.com/pkg/errors"

	"github.com/donglei1234/platform/services/auth/tools"
	"github.com/donglei1234/platform/services/common/nosql/document"
	errors2 "github.com/donglei1234/platform/services/common/nosql/errors"
)

type userAuthentication struct {
	ProfileId string `json:"profileId"`
	version   document.Version
}

func (d *Document) LoadOrCreateUserAuth(appId, name string) (*userAuthentication, error) {
	uid, err := tools.GenerateUUID()
	if err != nil {
		return nil, err
	}
	auth := &userAuthentication{
		ProfileId: uid,
	}

	if key, err := userAuthenticationKey(appId, name); err != nil {
		return nil, err
	} else if version, err := d.Get(key, document.WithDestination(auth)); errors.Cause(err) != errors2.ErrKeyNotFound {
		if err == nil {
			return auth, nil
		} else {
			return nil, err
		}
	} else if version, err = d.Set(key, document.WithSource(auth), document.WithAnyVersion()); err != nil {
		return nil, err
	} else {
		auth.version = version
	}
	return auth, nil
}

func (d *Document) LoadOrBindUserAuth(appId, platformType, platformUid, profileId string) (*userAuthentication, error) {
	auth := &userAuthentication{
		ProfileId: profileId,
	}

	if key, err := userPlatformAuthenticationKey(appId, platformType, platformUid); err != nil {
		return nil, err
	} else if version, err := d.Get(key, document.WithDestination(auth)); errors.Cause(err) != errors2.ErrKeyNotFound {
		if err == nil {
			auth.version = version
			return auth, nil
		} else {
			return nil, err
		}
	} else if version, err = d.Set(key, document.WithSource(auth), document.WithAnyVersion()); err != nil {
		return nil, err
	} else {
		auth.version = version
	}
	return auth, nil
}
