package db

import (
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/nosql/errors"
	pb "github.com/donglei1234/platform/services/proto/gen/mail/api"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	MailPublicKeyPrefix = "public"
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

func (db *Database) AddMails(target string, mails ...*pb.Mail) error {
	if len(mails) == 0 {
		return nil
	}
	key, err := makeMailKey(target)
	if err != nil {
		return err
	}
	uMails := make(map[string]interface{}, 0)
	for _, v := range mails {
		fieldKey, err := makeFieldMailKey(v.Id)
		if err != nil {
			return err
		}
		if msg, err := jsonx.Marshal(v); err != nil {
			return err
		} else {
			uMails[fieldKey.String()] = msg
		}
	}
	if res := db.HMSet(key.String(), uMails); res.Err() != nil {
		return res.Err()
	}
	return nil
}
func (db *Database) getPublicIndex(id string) (int32, error) {
	key, err := makeMailPublicIndexKey(id)
	if err != nil {
		return 0, err
	}
	if res := db.Get(key.String()); res.Err() != nil {
		if errors.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	} else {
		v, err := strconv.Atoi(res.Val())
		if err != nil {
			return 0, err
		}
		return int32(v), nil
	}
}

func (db *Database) savePublicIndex(id string, index int32) error {
	key, err := makeMailPublicIndexKey(id)
	if err != nil {
		return err
	}
	if res := db.Set(key.String(), index, -1); res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (db *Database) GetSelfMails(id string) map[int64]*pb.Mail {
	key, err := makeMailKey(id)
	if err != nil {
		return nil
	}
	if data := db.HGetAll(key.String()); data.Err() != nil {
		if errors.IsNotFound(data.Err()) {
			return make(map[int64]*pb.Mail)
		}
		return nil
	} else {
		res := make(map[int64]*pb.Mail, 0)
		for _, v := range data.Val() {
			m := &pb.Mail{}
			if err := jsonx.Unmarshal([]byte(v), m); err != nil {
				db.logger.Error("unmarshal error", zap.Error(err))
				continue
			}
			if m.Status == pb.MailStatus_DELETED {
				if err := db.DelMails(id, m.Id); err != nil {
					db.logger.Error("del mail error", zap.Error(err))
					continue
				}
			}
			db.checkAndMarkExpired(m)
			res[m.Id] = m
		}
		return res
	}
}

func (db *Database) checkAndMarkExpired(mail *pb.Mail) {
	if mail.Status == pb.MailStatus_DELETED {
		return
	}
	if mail.Expire > 0 && mail.Expire < time.Now().Unix() {
		mail.Status = pb.MailStatus_DELETED
	}
	return
}

func (db *Database) DelMails(profileId string, id int64) error {
	key, err := makeMailKey(profileId)
	if err != nil {
		return err
	}

	filedKey, err := makeFieldMailKey(id)
	if err != nil {
		return err
	}
	if res := db.HDel(key.String(), filedKey.String()); res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (db *Database) MergePublicAndPrivateMail(uid string) error {
	index, err := db.getPublicIndex(uid)
	if err != nil {
		return err
	}
	pMails, err := db.getPublicLeftMails(index)
	if err != nil {
		return err
	}
	if len(pMails) == 0 {
		return nil
	}
	if err := db.AddMails(uid, pMails...); err != nil {
		return err
	}
	rIndex := index + int32(len(pMails))
	if err := db.savePublicIndex(uid, rIndex); err != nil {
		return err
	}
	//TODO check mail num max and remove mail
	return nil
}
func (db *Database) UpdateAllStatus(profileId string, status pb.MailStatus) ([]*pb.Mail, error) {
	mAll := db.GetSelfMails(profileId)

	updates := make([]*pb.Mail, 0)
	for _, v := range mAll {
		if db.checkAndUpdateStatus(v, status) {
			updates = append(updates, v)
		}
	}
	err := db.AddMails(profileId, updates...)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (db *Database) PushMailToPublic(mail *pb.Mail) error {
	key, err := makeMailKey(MailPublicKeyPrefix)
	if err != nil {
		return err
	}
	m, err := jsonx.Marshal(mail)
	if err != nil {
		return err
	}
	if res := db.RPush(key.String(), m); res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (db *Database) getPublicLeftMails(start int32) ([]*pb.Mail, error) {
	key, err := makeMailKey(MailPublicKeyPrefix)
	if err != nil {
		return nil, err
	}
	if data := db.LRange(key.String(), int64(start), -1); data.Err() != nil {
		return nil, data.Err()
	} else {
		res := make([]*pb.Mail, 0)
		for _, v := range data.Val() {
			m := &pb.Mail{}
			if err := jsonx.Unmarshal([]byte(v), m); err != nil {
				db.logger.Error("unmarshal error", zap.Error(err))
				continue
			}
			res = append(res, m)
		}
		return res, nil
	}
}

func (db *Database) checkAndUpdateStatus(mail *pb.Mail, status pb.MailStatus) bool {
	if mail.Status >= status {
		return false
	}
	if status == pb.MailStatus_DELETED &&
		mail.Status != pb.MailStatus_REWARDED &&
		len(mail.Rewards) > 0 {
		return false
	}

	mail.Status = status
	return true
}

func (db *Database) getSelfAllMails(profileId string) ([]*pb.Mail, error) {
	key, err := makeMailKey(profileId)
	if err != nil {
		return nil, err
	}
	if data := db.LRange(key.String(), 0, -1); data.Err() != nil {
		return nil, data.Err()
	} else {
		res := make([]*pb.Mail, 0)
		for k, v := range data.Val() {
			m := &pb.Mail{}
			if err := jsonx.Unmarshal([]byte(v), m); err != nil {
				db.logger.Error("unmarshal error", zap.Error(err))
				continue
			}
			m.Id = int64(k)
			res = append(res, m)
		}
		return res, nil
	}

}

func (db *Database) UpdateOneStatus(profileId string, id int64, status pb.MailStatus) (*pb.Mail, error) {
	if id == 0 {
		return nil, errors.ErrInternal
	}
	key, err := makeMailKey(profileId)
	if err != nil {
		return nil, err
	}

	filedKey, err := makeFieldMailKey(id)
	if err != nil {
		return nil, err
	}
	mail := &pb.Mail{}
	if data := db.HGet(key.String(), filedKey.String()); data.Err() != nil {
		return nil, data.Err()
	} else {
		if err := jsonx.Unmarshal([]byte(data.Val()), mail); err != nil {
			return nil, err
		}
	}
	if !db.checkAndUpdateStatus(mail, status) {
		return nil, errors.ErrInternal
	}
	if data, err := jsonx.Marshal(mail); err != nil {
		return nil, err
	} else {
		if res := db.HSet(key.String(), filedKey.String(), data); res.Err() != nil {
			return nil, res.Err()
		} else {
			return mail, nil
		}
	}
}
