package db

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Database struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

const (
	Leaderboard   = ":leaderboard:"
	Option        = ":option"
	Reset         = ":reset"
	MethodType    = "MethodType"
	OrderType     = "OrderType"
	ResetTime     = "ResetTime"
	UpdateTime    = "UpdateTime"
	LatestVersion = "LatestVersion"
)

const (
	SUM    = "SUM"    // 将分数累积到总分
	BETTER = "BETTER" // 取历史最高分
	LAST   = "LAST"   // 取最新分数
)

const (
	DESCENDING = "DESCENDING" // high to low
	ASCENDING  = "ASCENDING"  // low to high
)

func OpenDatabase(l *zap.Logger, addr string, pwd string) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	go func() {
		l.Info("watching expired keys")
		watcher := client.PSubscribe("__keyspace@0__:*:leaderboard:*:reset").Channel()
		for {
			msg := <-watcher
			if msg.Payload == "expired" {
				keys := strings.Split(msg.Channel, ":")
				appId := keys[1]
				listId := keys[3]
				key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
				if exists := client.Exists(key).Val(); exists == 0 {
					panic(ErrListNotExists)
				}
				res, err := client.HGetAll(key).Result()
				if err != nil {
					panic(err)
				}
				resetTime, err := strconv.Atoi(res[ResetTime])
				if err != nil {
					panic(err)
				}
				version, err := strconv.Atoi(res[LatestVersion])
				if err != nil {
					panic(err)
				}
				if err := client.HSet(key, LatestVersion, version+1).Err(); err != nil {
					panic(err)
				}
				key = MakeLeaderboardKey(appId, Leaderboard, listId, Reset)
				if err := client.Set(key, 1, time.Duration(resetTime)*time.Second).Err(); err != nil {
					panic(err)
				}
				l.Debug("reset: " + msg.Channel)
			}
		}
	}()
	return &Database{
		l,
		client,
	}
}

func MakeLeaderboardKey(appId, leaderboard, listId, keyType string) string {
	return appId + leaderboard + listId + keyType
}

func (d *Database) NewLeaderboard(appId, listId, method, order string, resetTime, updateTime int32) error {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 1 {
		return ErrDuplicateListName
	}
	option := make(map[string]interface{})
	option[MethodType] = method
	option[OrderType] = order
	option[ResetTime] = resetTime
	option[UpdateTime] = updateTime
	option[LatestVersion] = 0
	if err := d.redisClient.HMSet(key, option).Err(); err != nil {
		return err
	}
	key = MakeLeaderboardKey(appId, Leaderboard, listId, Reset)
	if err := d.redisClient.Set(key, 1, time.Duration(resetTime)*time.Second).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetLeaderboardSize(appId, listId string) (int32, error) {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return 0, ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return 0, err
	}
	key = MakeLeaderboardKey(appId, Leaderboard, listId, ":"+strconv.Itoa(version))
	res, err := d.redisClient.ZCard(key).Result()
	return int32(res), err
}

func (d *Database) ResetLeaderboard(appId, listId string) error {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return err
	}
	if err := d.redisClient.HSet(key, LatestVersion, version+1).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateScore(appId, listId, id string, score float64) error {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return err
	}
	method := d.redisClient.HGet(key, MethodType).Val()
	key = MakeLeaderboardKey(appId, Leaderboard, listId, ":"+strconv.Itoa(version))
	rank := d.redisClient.ZRevRank(key, id).Val()
	obj := redis.Z{
		Score:  0,
		Member: id,
	}
	if rank == 0 {
		tem := d.redisClient.ZRevRange(key, 0, 0).Val()
		if len(tem) == 0 || tem[0] != id {
		} else {
			obj.Score = d.redisClient.ZScore(key, id).Val()
		}
	} else {
		obj.Score = d.redisClient.ZScore(key, id).Val()
	}
	switch method {
	case SUM:
		obj.Score += score
	case BETTER:
		if obj.Score < score {
			obj.Score = score
		}
	case LAST:
		obj.Score = score
	}
	if err := d.redisClient.ZAdd(key, obj).Err(); err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteMember(appId, listId, id string) error {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return err
	}
	key = MakeLeaderboardKey(appId, Leaderboard, listId, ":"+strconv.Itoa(version))
	if ok := d.redisClient.ZRem(key, id).Val(); ok == 0 {
		return ErrMemberNotExists
	}
	return nil
}
func (d *Database) GetRankById(appId, listId, id string) (int32, error) {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return 0, ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return 0, err
	}
	order := d.redisClient.HGet(key, OrderType).Val()
	key = MakeLeaderboardKey(appId, Leaderboard, listId, ":"+strconv.Itoa(version))
	rank := d.redisClient.ZRevRank(key, id).Val()
	if rank == 0 {
		tem := d.redisClient.ZRevRange(key, 0, 0).Val()
		if len(tem) == 0 || tem[0] != id {
			return 0, ErrMemberNotExists
		}
	}
	if order == ASCENDING {
		rank = d.redisClient.ZRank(key, id).Val()
	}
	return int32(rank), nil
}

func (d *Database) GetRankFromMToN(appId, listId string, m, n int64) ([]redis.Z, error) {
	key := MakeLeaderboardKey(appId, Leaderboard, listId, Option)
	if exists := d.redisClient.Exists(key).Val(); exists == 0 {
		return nil, ErrListNotExists
	}
	version, err := strconv.Atoi(d.redisClient.HGet(key, LatestVersion).Val())
	if err != nil {
		return nil, err
	}
	order := d.redisClient.HGet(key, OrderType).Val()
	key = MakeLeaderboardKey(appId, Leaderboard, listId, ":"+strconv.Itoa(version))
	if order == DESCENDING {
		return d.redisClient.ZRevRangeWithScores(key, m, n).Result()
	} else {
		return d.redisClient.ZRangeWithScores(key, m, n).Result()
	}
}
