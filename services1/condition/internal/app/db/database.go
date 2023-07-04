package db

import (
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/nosql/errors"
	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/options"
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"github.com/donglei1234/platform/services/condition/internal/app/db/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
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

func (db *Database) AddCondition(profileId string, adds ...*pb.Condition) error {
	key, err := makeConditionKey(profileId)
	if err != nil {
		return err
	}
	conds, err := db.getConditions(key, adds...)
	if err != nil {
		return err
	}

	for _, v := range adds {
		fieldKey, err := makeConditionFieldKey(v)
		if err != nil {
			return err
		}
		_, ok := conds[fieldKey.String()]
		if !ok {
			conds[fieldKey.String()] = model.NewCondition()
		}
		conds[fieldKey.String()].AddCondition(v)
	}

	dSet := make(map[string]interface{})
	for k, v := range conds {
		dSet[k] = v
	}
	err = db.HMSet(key, options.WithMultipleSource(dSet))
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) getConditions(key keys.Key, conds ...*pb.Condition) (map[string]*model.Condition, error) {
	if len(conds) == 0 {
		return make(map[string]*model.Condition), nil
	}
	addsType := make(map[keys.Key]struct{}, 0)
	for _, cond := range conds {
		fieldKey, err := makeConditionFieldKey(cond)
		if err != nil {
			return nil, err
		}
		addsType[fieldKey] = struct{}{}
	}

	gKeys := lo.Keys[keys.Key, struct{}](addsType)
	conditionsGet := make(map[string]interface{}, len(gKeys))

	err := db.HMGet(key, gKeys, options.WithDestinationMap(conditionsGet))
	if err != nil {
		if err != errors.ErrKeyNotFound {
			return nil, err
		}
	}
	result := make(map[string]*model.Condition)
	for k, v := range conditionsGet {
		cond := model.NewCondition()
		err = jsonx.Unmarshal([]byte(v.(string)), cond)
		if err != nil {
			db.logger.Error("unmarshal condition error", zap.Error(err))
			continue
		}
		result[k] = cond
	}
	return result, nil
}

func (db *Database) UpdateCondition(profileId string, updates ...*pb.Condition) error {
	key, err := makeConditionKey(profileId)
	if err != nil {
		return err
	}

	condsData, err := db.getConditions(key, updates...)
	if err != nil {
		return err
	}
	if len(condsData) == 0 {
		return nil
	}

	for _, v := range updates {
		fieldKey, err := makeConditionFieldKey(v)
		if err != nil {
			return err
		}
		data, ok := condsData[fieldKey.String()]
		if !ok {
			continue
		}
		data.UpdateCondition(v)
	}
	dSet := make(map[string]interface{})
	for k, v := range condsData {
		dSet[k] = v
	}
	return db.HMSet(key, options.WithMultipleSource(dSet))
}

func (db *Database) Clear(profileId string) error {
	key, err := makeConditionKey(profileId)
	if err != nil {
		return err

	}
	return db.Del(key)
}

func (db *Database) GetAndDeleteFinishedConditions(profileId string, condition *pb.Condition) ([]*pb.Condition, error) {
	key, err := makeConditionKey(profileId)
	if err != nil {
		return nil, err
	}

	fieldKey, err := makeConditionFieldKey(condition)
	if err != nil {
		return nil, err
	}
	cond := model.NewCondition()
	err = db.HGet(key, fieldKey, options.WithDestination(cond))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	res := cond.GetConditionLst()
	cond.RemoveFinishedCondition()
	if cond.IsEmpty() {
		err := db.HDel(key, fieldKey)
		if err != nil {
			return nil, err
		}
	} else {
		err = db.HSet(key, fieldKey, options.WithSource(cond))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (db *Database) DeleteCondition(profileId string, conditions ...*pb.Condition) error {
	key, err := makeConditionKey(profileId)
	if err != nil {
		return err
	}

	conds, err := db.getConditions(key, conditions...)
	if err != nil {
		return err
	}

	for _, v := range conditions {
		fieldKey, err := makeConditionFieldKey(v)
		if err != nil {
			return err
		}
		data, ok := conds[fieldKey.String()]
		if !ok {
			continue
		}
		data.RemoveCondition(v)
	}
	dSet := make(map[string]interface{})
	for k, v := range conds {
		dSet[k] = v
	}
	return db.HMSet(key, options.WithMultipleSource(dSet))
}
