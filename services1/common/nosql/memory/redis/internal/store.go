package internal

import (
	"reflect"
	"strconv"

	"github.com/cskr/pubsub"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/syncmap"

	"github.com/donglei1234/platform/services/common/jsonx"
	nerr "github.com/donglei1234/platform/services/common/nosql/errors"
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/options"
)

type MemoryStore struct {
	logger *zap.Logger
	name   string
	client *redis.Client
	subs   syncmap.Map
	pubsub *pubsub.PubSub
}

func NewMemoryStore(name string, client *redis.Client, l *zap.Logger) (*MemoryStore, error) {
	return &MemoryStore{
		name:   name,
		client: client,
		logger: l,
		pubsub: pubsub.New(0),
	}, nil
}

func (d *MemoryStore) Name() string {
	return d.name
}

func (d *MemoryStore) Incur(key keys.Key) (int64, error) {
	if score, err := d.client.Incr(key.String()).Result(); err != nil {
		return 0, err
	} else {
		return score, nil
	}
}

func (d *MemoryStore) Set(key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if data, err := jsonx.Marshal(o.Source); err != nil {
		return err
	} else if value, err := d.client.Set(key.String(), data, o.TTL).Result(); err != nil {
		return err
	} else if o.Destination != nil {
		return jsonx.Unmarshal([]byte(value), o.Destination)
	}
	return nil
}

func (d *MemoryStore) MSet(opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if len(o.SourceList) <= 0 {
		return errors.New("source list is empty")
	} else if value, err := d.client.MSet(o.SourceList).Result(); err != nil {
		return err
	} else if o.Destination != nil {
		return jsonx.Unmarshal([]byte(value), o.Destination)
	}
	return nil
}

func (d *MemoryStore) SetNX(key keys.Key, opts ...options.Option) (bool, error) {
	if o, err := options.NewOptions(opts...); err != nil {
		return false, err
	} else if data, err := jsonx.Marshal(o.Source); err != nil {
		return false, err
	} else if isOk, err := d.client.SetNX(key.String(), data, o.TTL).Result(); err != nil {
		return false, err
	} else {
		return isOk, nil
	}
}

func (d *MemoryStore) Lock(key keys.Key) error {
	////TODO 随机生成token 存储
	//token := ""
	//if isOk, err := d.client.SetNX(key.String(), token, LockTTL).Result(); err != nil {
	//	return err
	//}else if
	return nil

}

func (d *MemoryStore) Get(key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if value, err := d.client.Get(key.String()).Result(); err != nil {
		return convertRedisError(err, "")
	} else if o.Destination == nil {
		return nerr.ErrDestIsNil
	} else {
		return jsonx.Unmarshal([]byte(value), o.Destination)
	}
}

func (d *MemoryStore) Del(key keys.Key) error {
	if _, err := d.client.Del(key.String()).Result(); err != nil {
		return err
	} else {
		return nil
	}
}

func (d *MemoryStore) Exists(key keys.Key) (bool, error) {
	if isExists, err := d.client.Exists(key.String()).Result(); err != nil {
		return false, err
	} else if isExists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *MemoryStore) HSet(path keys.Key, key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if data, err := jsonx.Marshal(o.Source); err != nil {
		return err
	} else if _, err := d.client.HSet(path.String(), key.String(), data).Result(); err != nil {
		return err
	}
	return nil
}

func (d *MemoryStore) HGet(path keys.Key, key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.Destination == nil {
		return nerr.ErrDestIsNil
	} else if value, err := d.client.HGet(path.String(), key.String()).Result(); err != nil {
		return convertRedisError(err, "")
	} else {
		return jsonx.Unmarshal([]byte(value), o.Destination)
	}
}

func (d *MemoryStore) HDel(path keys.Key, keys ...keys.Key) error {
	keyStrings := make([]string, len(keys))
	for k, v := range keys {
		keyStrings[k] = v.String()
	}

	if _, err := d.client.HDel(path.String(), keyStrings...).Result(); err != nil {
		return err
	}
	return nil
}

func (d *MemoryStore) HLen(path keys.Key) (int64, error) {
	if length, err := d.client.HLen(path.String()).Result(); err != nil {
		return 0, err
	} else {
		return length, nil
	}
}

func (d *MemoryStore) HMSet(path keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else {
		dSets := make(map[string]interface{})
		for k, v := range o.Sources {
			if data, err := jsonx.Marshal(v); err != nil {
				return err
			} else {
				dSets[k] = data
			}
		}

		if _, err := d.client.HMSet(path.String(), dSets).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (d *MemoryStore) HGetAll(path keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if values, err := d.client.HGetAll(path.String()).Result(); err != nil {
		return err
	} else {
		if o.DestinationMap != nil {
			for k, v := range values {
				o.DestinationMap[k] = v
			}
		}
	}
	return nil
}

func (d *MemoryStore) HMGet(path keys.Key, keys []keys.Key, opts ...options.Option) error {
	keyStrings := make([]string, len(keys))
	for k, v := range keys {
		keyStrings[k] = v.String()
	}

	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if values, err := d.client.HMGet(path.String(), keyStrings...).Result(); err != nil {
		return err
	} else {
		if o.DestinationMap != nil {
			for k, v := range values {
				if v == nil {
					continue
				}
				o.DestinationMap[keys[k].String()] = v
			}
		}
	}
	return nil
}

func (d *MemoryStore) makeRedisZ(in map[string]interface{}) ([]redis.Z, error) {
	items := make([]redis.Z, len(in))
	index := 0
	for k, v := range in {
		v, ok := v.(float64)
		if !ok {
			return nil, nerr.ErrNotFloat64
		}

		items[index].Member = k
		items[index].Score = v

		index++
	}
	return items, nil
}

// ZAdd: https://redis.io/commands/zadd
func (d *MemoryStore) ZAdd(path keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.Sources == nil || len(o.Sources) <= 0 {
		return nerr.ErrDestIsNil
	} else if redisZs, err := d.makeRedisZ(o.Sources); err != nil {
		return err
	} else if _, err := d.client.ZAdd(path.String(), redisZs...).Result(); err != nil {
		return err
	}
	return nil
}

// ZScore: https://redis.io/commands/zscore
func (d *MemoryStore) ZScore(path keys.Key, key keys.Key) (float64, error) {
	if score, err := d.client.ZScore(path.String(), key.String()).Result(); err != nil {
		return 0, err
	} else {
		return score, nil
	}
}

// ZRevRank: https: //redis.io/commands/zrevrank
func (d *MemoryStore) ZRevRank(path keys.Key, key keys.Key) (int64, error) {
	if rank, err := d.client.ZRevRank(path.String(), key.String()).Result(); err != nil {
		return 0, err
	} else {
		return rank, nil
	}
}

// ZRank: https://redis.io/commands/zrank
func (d *MemoryStore) ZRank(path keys.Key, key keys.Key) (int64, error) {
	if rank, err := d.client.ZRank(path.String(), key.String()).Result(); err != nil {
		return 0, err
	} else {
		return rank, nil
	}
}

func (d *MemoryStore) ZRevRange(path keys.Key, destOpt options.Option, opts ...options.ZRangeOption) error {
	//if o, e := options.NewZRangeOption(opts...); e != nil {
	//	return e
	//} else if do, e := options.NewOptions(destOpt); e != nil {
	//	return e
	//} else if do.DestinationList == nil || int64(len(do.DestinationList)) <= (o.End-o.Start) {
	//	return nerr.ErrDestIsNil
	//} else if rankList, err := d.client.ZRevRange(path.String(), o.Start, o.End).Result(); err != nil {
	//	return err
	//} else {
	//	for k, v := range rankList {
	//		do.DestinationList[k] = v
	//	}
	//}
	return nil
}

func (d *MemoryStore) makeBytes(in []interface{}) ([]interface{}, error) {
	items := make([]interface{}, len(in))
	for index, v := range in {
		bytes, err := jsonx.Marshal(v)
		if err != nil {
			return nil, err
		}
		items[index] = bytes
	}
	return items, nil
}

func (d *MemoryStore) LPush(key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.SourceList == nil || len(o.SourceList) <= 0 {
		return nerr.ErrDestIsNil
	} else if list, err := d.makeBytes(o.SourceList); err != nil {
		return err
	} else if _, err = d.client.LPush(key.String(), list...).Result(); err != nil {
		return err
	}
	return nil
}

func (d *MemoryStore) LRange(key keys.Key, ro options.ZRangeOption, opts ...options.Option) error {
	if o, err := options.NewZRangeOption(ro); err != nil {
		return err
	} else if ops, err := options.NewOptions(opts...); err != nil {
		return nerr.ErrDestIsNil
	} else if list, err := d.client.LRange(key.String(), o.Start, o.End).Result(); err != nil {
		return err
	} else {
		for k, v := range list {
			ops.DestinationMap[strconv.Itoa(k)] = v
		}
	}
	return nil
}

func (d *MemoryStore) RPush(key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.SourceList == nil || len(o.SourceList) <= 0 {
		return nerr.ErrDestIsNil
	} else if list, err := d.makeBytes(o.SourceList); err != nil {
		return err
	} else if _, err = d.client.RPush(key.String(), list...).Result(); err != nil {
		return err
	}
	return nil
}

func (d *MemoryStore) LPop(key keys.Key, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.Destination == nil {
		return nerr.ErrDestIsNil
	} else if value, err := d.client.LPop(key.String()).Result(); err != nil {
		return err
	} else {
		return jsonx.Unmarshal([]byte(value), o.Destination)
	}
}

func (d *MemoryStore) LLen(key keys.Key) (int64, error) {
	if length, err := d.client.LLen(key.String()).Result(); err != nil {
		return 0, err
	} else {
		return length, nil
	}
}

func (d *MemoryStore) LSet(key keys.Key, index int64, opts ...options.Option) error {
	if o, err := options.NewOptions(opts...); err != nil {
		return err
	} else if o.Source == nil {
		return nerr.ErrDestIsNil
	} else if value, err := jsonx.Marshal(o.Source); err != nil {
		return err
	} else if _, err = d.client.LSet(key.String(), index, value).Result(); err != nil {
		return err
	}
	return nil
}

// Scan https://redis.io/commands/scan
func (d *MemoryStore) Scan(path keys.Key, destOpt options.Option, scanOpt ...options.ScanOption) (cursor int, err error) {
	// make sure we have a destination
	if destOpt == nil {
		err = nerr.ErrDestIsNil
		return
	}
	opts := &options.Options{}
	if e := destOpt(opts); e != nil {
		err = e
	} else if o, e := options.NewScanOptions(scanOpt...); e != nil {
		err = e
	} else if opts.Destination == nil {
		return 0, nerr.ErrDestIsNil
	} else if len(o.Query) == 0 {
		return 0, errors.Wrap(nerr.ErrNoScanType, "scan type unset")
	} else {
		dest := reflect.ValueOf(opts.Destination).Elem()
		k := dest.Kind()
		if k != reflect.Map {
			err = errors.Wrap(nerr.ErrInternal, "unsupported destination type")
			return
		}

		res, cur, e := d.scan(path, o.Query[0], o.Offset, o.Limit)
		if e != nil {
			err = e
			return
		}
		cursor = int(cur)

		if len(res) <= 0 {
			return
		}

		for i := 0; i < len(res); i += 2 {
			key := res[i]
			value := []byte(res[i+1])
			t := dest.Type().Elem()
			buf := reflect.New(t).Interface()
			if err = jsonx.Unmarshal(value, &buf); err != nil {
				return
			} else {
				k := reflect.ValueOf(key)
				v := reflect.ValueOf(buf).Elem()
				dest.SetMapIndex(k, v)
			}
		}
	}
	return
}

func (d *MemoryStore) scan(path keys.Key, query options.ScanQuery, offset int, limit int) ([]string, uint64, error) {
	if query.ScanType != options.ScanTypeRegex {
		return nil, 0, errors.Wrap(nerr.ErrNotImplemented, "unsupported scan type")
	}
	return d.client.HScan(path.String(), uint64(offset), query.Regex, int64(limit)).Result()
}

func (d *MemoryStore) Remove(keys ...keys.Key) error {
	ks := make([]string, len(keys))
	for k, v := range keys {
		ks[k] = v.String()
	}
	_, err := d.client.Del(ks...).Result()
	if err != nil {
		return err
	}
	return nil
}

const PSubscribeChannel = "__keyspace@*__:"

func (d *MemoryStore) Sub(key keys.Key) (chan interface{}, error) {
	sub := d.client.PSubscribe(PSubscribeChannel + key.String())
	_, isLoad := d.subs.LoadOrStore(key.String(), sub)
	if !isLoad {
		ch := sub.Channel()
		go func(key string) {
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						d.pubsub.Close(key)
						return
					}
					d.pubsub.Pub(keys.ChangesType(msg.Payload), key)
				}
			}
		}(key.String())
	}

	c := d.pubsub.Sub(key.String())
	return c, nil
}

func (d *MemoryStore) UnSub(key keys.Key) error {
	value, ok := d.subs.LoadAndDelete(key.String())
	if !ok {
		return nerr.ErrKeyNotFound
	}
	sub := value.(*redis.PubSub)
	err := sub.PUnsubscribe(PSubscribeChannel + key.String())
	if err != nil {
		return err
	}
	err = sub.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *MemoryStore) CloseSub(ch chan interface{}) {
	d.pubsub.Unsub(ch)
}

func (d *MemoryStore) ShutDown() {
	d.pubsub.Shutdown()
}

func convertRedisError(e error, key string) error {
	if e == nil {
		return nil
	}

	switch e {
	case redis.Nil:
		if key == "" {
			return nerr.ErrKeyNotFound
		} else {
			return errors.Wrap(nerr.ErrKeyNotFound, key)
		}
	default:
		return errors.Wrap(nerr.ErrDriverFailure, e.Error())
	}
}
