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
	NotificationUser     = ":notification:user:"
	NotificationDeviceId = ":notification:deviceId:"
	NotificationTopic    = ":notification:topic"
	DeviceId             = "deviceId"
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

func MakeNotificationKey(appId, notification, profileId string) string {
	return appId + notification + profileId
}

func (d *Database) GetUserInfoByDeviceId(profileId, appId, deviceId string) (map[string]string, error) {
	var deviceToken, region, endPointArn string
	userInfo := make(map[string]string)
	deviceKey := MakeNotificationKey(appId, NotificationDeviceId, deviceId)
	if exist := d.redisClient.Exists(deviceKey).Val(); exist == 0 {
		return userInfo, ErrUserKeyNotExist
	}
	deviceToken = d.redisClient.HGet(deviceKey, "deviceToken").Val()
	region = d.redisClient.HGet(deviceKey, "region").Val()
	endPointArn = d.redisClient.HGet(deviceKey, "endPointArn").Val()
	userInfo["deviceToken"] = deviceToken
	userInfo["region"] = region
	userInfo["endPointArn"] = endPointArn
	return userInfo, nil
}

func (d *Database) IsDeviceKeyExist(appId, deviceId string) bool {
	deviceKey := MakeNotificationKey(appId, NotificationDeviceId, deviceId)
	if exist := d.redisClient.Exists(deviceKey).Val(); exist == 0 {
		return false
	}
	return true
}

func (d *Database) SetDeviceId(profileId, appId, deviceId string) error {
	userKey := MakeNotificationKey(appId, NotificationUser, profileId)
	if err := d.redisClient.HSet(userKey, DeviceId, deviceId).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) DelUserInfoByProfileId(profileId, appId, deviceId string) error {
	userKey := MakeNotificationKey(appId, NotificationUser, profileId)
	deviceIdKey := MakeNotificationKey(appId, NotificationDeviceId, deviceId)
	if err := d.redisClient.Del(userKey).Err(); err != nil {
		return err
	}
	if err := d.redisClient.Del(deviceIdKey).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) SetDeviceIdAndUserInfo(profileId, appId, deviceId, deviceToken, region, deviceType, endPointArn string) error {
	userKey := MakeNotificationKey(appId, NotificationUser, profileId)
	deviceIdKey := MakeNotificationKey(appId, NotificationDeviceId, deviceId)
	if err := d.redisClient.HSet(userKey, DeviceId, deviceId).Err(); err != nil {
		return err
	}
	userInfo := map[string]interface{}{
		"deviceToken": deviceToken,
		"region":      region,
		"endPointArn": endPointArn,
		"deviceType":  deviceType,
	}
	if err := d.redisClient.HMSet(deviceIdKey, userInfo).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetDeviceIdByProfileId(profileId, appId string) (deviceId string, err error) {
	userKey := MakeNotificationKey(appId, NotificationUser, profileId)
	if exist := d.redisClient.Exists(userKey).Val(); exist == 0 {
		return deviceId, ErrUserKeyNotExist
	}
	deviceId = d.redisClient.HGet(userKey, DeviceId).Val()
	return deviceId, nil
}

func (d *Database) SetTopicArn(appId, topicName, topicArn string) error {
	topicKey := MakeNotificationKey(appId, NotificationTopic, "")
	if err := d.redisClient.HSet(topicKey, topicName, topicArn).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetTopicArn(appId, topicName string) (string, error) {
	topicKey := MakeNotificationKey(appId, NotificationTopic, "")
	if exist := d.redisClient.Exists(topicKey).Val(); exist == 0 {
		return "", ErrTopicKeyNotExist
	}
	topicArn := d.redisClient.HGet(topicKey, topicName).Val()
	return topicArn, nil
}

func (d *Database) DelTopicArn(appId, TopicName string) error {
	topicKey := MakeNotificationKey(appId, NotificationTopic, "")
	if exist := d.redisClient.Exists(topicKey, TopicName).Val(); exist == 0 {
		return ErrTopicKeyNotExist
	}
	if err := d.redisClient.HDel(topicKey, TopicName).Err(); err != nil {
		return err
	}
	return nil
}
