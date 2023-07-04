package db

import (
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/nosql/errors"
	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/options"
	pb "github.com/donglei1234/platform/services/proto/gen/gm/api"
	"go.uber.org/zap"
)

const (
	BanProfilePrefix = "profiles"
)

type Database struct {
	memory.MemoryStore
	logger *zap.Logger
}

func OpenDatabase(l *zap.Logger, mm memory.MemoryStore) *Database {
	return &Database{
		mm,
		l,
	}
}

func (db *Database) SetProfileBanStatus(info *pb.BanInfo, status pb.SetProfilesBanStatusRequest_BanStatus) error {
	path, err := makeProfileBanPath(BanProfilePrefix)
	if err != nil {
		return err
	}
	key, err := makeProfileBanKey(info.ProfileId)
	if err != nil {
		return err
	}
	if status == pb.SetProfilesBanStatusRequest_BAN_STATUS_UNBAN {
		err = db.HDel(path, key)
		if err != nil {
			return err
		}
	} else {
		err = db.HSet(path, key, options.WithSource(info))
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) GetAllProfilesBan() ([]*pb.BanInfo, error) {
	return db.GetProfilesBanInfo(make([]string, 0))
}

func (db *Database) GetProfilesBanInfo(profileIds []string) ([]*pb.BanInfo, error) {
	path, err := makeProfileBanPath(BanProfilePrefix)
	if err != nil {
		return nil, err
	}
	if len(profileIds) == 0 {
		return db.getAllBanInfos(path)
	} else {
		return db.getProfileBanInfo(path, profileIds...)
	}
}

func (db *Database) getProfileBanInfo(path keys.Key, profileIds ...string) ([]*pb.BanInfo, error) {
	gKeys := make([]keys.Key, len(profileIds))
	for i, v := range profileIds {
		k, err := makeProfileBanKey(v)
		if err != nil {
			return nil, err
		}
		gKeys[i] = k
	}
	dests := make(map[string]interface{}, len(gKeys))
	err := db.HMGet(path, gKeys, options.WithDestinationMap(dests))
	if err != nil {
		if err != errors.ErrKeyNotFound {
			return nil, err
		}
	}
	infos := make([]*pb.BanInfo, 0)
	for _, v := range dests {
		m := &pb.BanInfo{}
		if err := jsonx.Unmarshal([]byte(v.(string)), m); err != nil {
			db.logger.Error("unmarshal error", zap.Error(err))
			continue
		}
		infos = append(infos, m)
	}
	return infos, nil
}

func (db *Database) getAllBanInfos(path keys.Key) ([]*pb.BanInfo, error) {
	dest := make(map[string]interface{})
	err := db.HGetAll(path, options.WithDestinationMap(dest))
	if err != nil {
		return nil, err
	}
	infos := make([]*pb.BanInfo, 0)
	for _, v := range dest {
		m := &pb.BanInfo{}
		if err := jsonx.Unmarshal([]byte(v.(string)), m); err != nil {
			db.logger.Error("unmarshal error", zap.Error(err))
			continue
		}
		infos = append(infos, m)
	}
	return infos, nil
}

func (db *Database) SetBulletin(info *pb.BulletinInfo) error {
	path, err := makeBulletinPath()
	if err != nil {
		return err
	}
	key, err := makeBulletinKey(info.Id)
	if err != nil {
		return err
	}
	err = db.HSet(path, key, options.WithSource(info))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetBulletins() ([]*pb.BulletinInfo, error) {
	path, err := makeBulletinPath()
	if err != nil {
		return nil, err
	}
	dest := make(map[string]interface{})
	err = db.HGetAll(path, options.WithDestinationMap(dest))
	if err != nil {
		return nil, err
	}
	infos := make([]*pb.BulletinInfo, 0)
	for _, v := range dest {
		m := &pb.BulletinInfo{}
		if err := jsonx.Unmarshal([]byte(v.(string)), m); err != nil {
			db.logger.Error("unmarshal error", zap.Error(err))
			continue
		}
		infos = append(infos, m)
	}
	return infos, nil
}
