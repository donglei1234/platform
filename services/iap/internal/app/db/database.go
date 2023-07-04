package db

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Database struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

const (
	IAPUser     = ":iap:user:"
	IAPDeviceId = ":iap:deviceId:"
	DeviceId    = "deviceId"
)

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

func (db *Database) InsertIAPTradingRecord(appId string, token string, msg string) (bool, error) {
	log := db.redisClient.HSet(appId, token, msg)
	status, err := log.Result()
	return status, err
}

func (db *Database) GetIAPTradingRecord(key, field string) (string, error) {
	return db.redisClient.HGet(key, field).Result()
}

func (db *Database) CheckToken(key, field string) (bool, error) {
	return db.redisClient.HExists(key, field).Result()
}

func (db *Database) UpdateIAPTradingRecord(appId string, token string, msg string) (bool, error) {
	log := db.redisClient.HSet(appId, token, msg)
	status, err := log.Result()
	return status, err
}
