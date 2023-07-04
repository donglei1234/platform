package db

import (
	"github.com/donglei1234/platform/services/common/jsonx"
	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

type Database struct {
	*redis.Client
	logger *zap.Logger
}

func OpenDatabase(l *zap.Logger, client *redis.Client) *Database {
	return &Database{
		client,
		l,
	}
}

var messageExpire = time.Hour * 24 * 30

func (db *Database) PushChatMessage(uid string, msgs ...*pb.ChatMessage_Message) error {
	if key, e := makeChatMsgKey(uid); e != nil {
		return e
	} else {
		sData := make([]interface{}, 0)
		for _, v := range msgs {
			if data, err := jsonx.Marshal(v); err != nil {
				return err
			} else {
				sData = append(sData, data)
			}
		}

		if res := db.RPush(key.String(), sData...); res.Err() != nil {
			return res.Err()
		} else {
			db.Expire(key.String(), messageExpire)
		}
	}
	return nil
}

func (db *Database) GetAndDeleteChatMessage(uid string) ([]*pb.ChatMessage_Message, error) {
	if key, e := makeChatMsgKey(uid); e != nil {
		return nil, e
	} else {
		if data, e := db.LRange(key.String(), 0, -1).Result(); e != nil {
			return nil, e
		} else {
			res := make([]*pb.ChatMessage_Message, 0)
			for _, v := range data {
				msg := &pb.ChatMessage_Message{}
				if e := jsonx.Unmarshal([]byte(v), msg); e != nil {
					return nil, e
				}
				res = append(res, msg)
			}
			db.Del(key.String())
			return res, nil
		}
	}
}
