package db

import (
	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

type Database struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

const (
	FileKeyPathMap       = ":s3:key:path:"
	GlobalKeyPathProfile = "global"
)

func MakeStorageKey(appId, storage, profileId string) string {
	return appId + storage + profileId
}

func OpenDatabase(l *zap.Logger, addr string, pwd string) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	return &Database{
		l,
		client,
	}
}

func (d *Database) SavePathKey(profileId, appId, fileKey string, key string) (bool, error) {
	if profileId == "" {
		storageKey := MakeStorageKey(appId, FileKeyPathMap, GlobalKeyPathProfile)
		return d.redisClient.HSet(storageKey, fileKey, key).Result()
	} else {
		storageKey := MakeStorageKey(appId, FileKeyPathMap, profileId)
		return d.redisClient.HSet(storageKey, fileKey, key).Result()
	}
}

func (d *Database) GetPathByKey(profileId, appId, key string) (string, error) {
	if profileId == "" {
		storageKey := MakeStorageKey(appId, FileKeyPathMap, GlobalKeyPathProfile)
		return d.redisClient.HGet(storageKey, key).Result()
	} else {
		storageKey := MakeStorageKey(appId, FileKeyPathMap, profileId)
		return d.redisClient.HGet(storageKey, key).Result()
	}
}
