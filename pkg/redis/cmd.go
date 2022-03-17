// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/imkuqin-zw/courier/pkg/utils/xmap"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type (
	Z         = redis.Z
	ZWithKey  = redis.ZWithKey
	ZStore    = redis.ZStore
	Pipeliner = redis.Pipeliner
)

func New() *Redis {
	return RawConfig("redis").Build()
}

func (r *Redis) Ping() error {
	return r.UniversalClient.Ping().Err()
}

func (r *Redis) TTL(key string) (time.Duration, error) {
	return r.UniversalClient.TTL(key).Result()
}

func (r *Redis) Type(key string) (string, error) {
	return r.UniversalClient.Type(key).Result()
}

func (r *Redis) Decr(key string) (int64, error) {
	return r.UniversalClient.Decr(key).Result()
}

func (r *Redis) DecrBy(key string, decrement int64) (int64, error) {
	return r.UniversalClient.DecrBy(key, decrement).Result()
}

func (r *Redis) GetString(key string) (string, error) {
	result, err := r.UniversalClient.Get(key).Result()
	if err != redis.Nil {
		return result, err
	}
	return result, redis.Nil
}

func (r *Redis) Del(keys ...string) (int64, error) {
	return r.UniversalClient.Del(keys...).Result()
}

func (r *Redis) Incr(key string) (int64, error) {
	return r.UniversalClient.Incr(key).Result()
}

func (r *Redis) IncrBy(key string, value int64) (int64, error) {
	return r.UniversalClient.IncrBy(key, value).Result()
}

func (r *Redis) IncrByFloat(key string, value float64) (float64, error) {
	return r.UniversalClient.IncrByFloat(key, value).Result()
}

func (r *Redis) MGet(keys ...string) ([]interface{}, error) {
	return r.UniversalClient.MGet(keys...).Result()
}

func (r *Redis) MSet(pairs ...interface{}) error {
	return r.UniversalClient.MSet(pairs...).Err()
}

func (r *Redis) MSetNX(pairs ...interface{}) (bool, error) {
	return r.UniversalClient.MSetNX(pairs...).Result()
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.UniversalClient.Set(key, value, expiration).Err()
}

func (r *Redis) SetBit(key string, offset int64, value int) (int64, error) {
	return r.UniversalClient.SetBit(key, offset, value).Result()
}

func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.UniversalClient.SetNX(key, value, expiration).Result()
}

func (r *Redis) SetRange(key string, offset int64, value string) (int64, error) {
	return r.UniversalClient.SetRange(key, offset, value).Result()
}

func (r *Redis) StrLen(key string) (int64, error) {
	return r.UniversalClient.StrLen(key).Result()
}

func (r *Redis) HDel(key string, fields ...string) (int64, error) {
	return r.UniversalClient.HDel(key, fields...).Result()
}

func (r *Redis) HExists(key, field string) (bool, error) {
	return r.UniversalClient.HExists(key, field).Result()
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return r.UniversalClient.HGetAll(key).Result()
}

func (r *Redis) HIncrBy(key, field string, incr int64) (int64, error) {
	return r.UniversalClient.HIncrBy(key, field, incr).Result()
}

func (r *Redis) HIncrByFloat(key, field string, incr float64) (float64, error) {
	return r.UniversalClient.HIncrByFloat(key, field, incr).Result()
}

func (r *Redis) HKeys(key string) ([]string, error) {
	return r.UniversalClient.HKeys(key).Result()
}

func (r *Redis) HLen(key string) (int64, error) {
	return r.UniversalClient.HLen(key).Result()
}

func (r *Redis) HMGet(key string, fields ...string) ([]interface{}, error) {
	return r.UniversalClient.HMGet(key).Result()
}

func (r *Redis) HMSet(key string, fields map[string]interface{}) error {
	return r.UniversalClient.HMSet(key, fields).Err()
}

func (r *Redis) HSet(key, field string, value interface{}) (bool, error) {
	return r.UniversalClient.HSet(key, field, value).Result()
}

func (r *Redis) HSetNX(key, field string, value interface{}) (bool, error) {
	return r.UniversalClient.HSetNX(key, field, value).Result()
}

func (r *Redis) HVals(key string) ([]string, error) {
	return r.UniversalClient.HVals(key).Result()
}

func (r *Redis) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.UniversalClient.BLPop(timeout, keys...).Result()
}

func (r *Redis) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.UniversalClient.BRPop(timeout, keys...).Result()
}

func (r *Redis) BRPopLPush(source, destination string, timeout time.Duration) (string, error) {
	return r.UniversalClient.BRPopLPush(source, destination, timeout).Result()
}

func (r *Redis) LIndex(key string, index int64) (string, error) {
	val, err := r.UniversalClient.LIndex(key, index).Result()
	if err != redis.Nil {
		return val, err
	}
	return val, redis.Nil
}

func (r *Redis) LInsert(key, op string, pivot, value interface{}) (int64, error) {
	return r.UniversalClient.LInsert(key, op, pivot, value).Result()
}

func (r *Redis) LInsertBefore(key string, pivot, value interface{}) (int64, error) {
	return r.UniversalClient.LInsertBefore(key, pivot, value).Result()
}

func (r *Redis) LInsertAfter(key string, pivot, value interface{}) (int64, error) {
	return r.UniversalClient.LInsertBefore(key, pivot, value).Result()
}

func (r *Redis) LLen(key string) (int64, error) {
	return r.UniversalClient.LLen(key).Result()
}

func (r *Redis) LPop(key string) (string, error) {
	val, err := r.UniversalClient.LPop(key).Result()
	if err != redis.Nil {
		return val, err
	}
	return val, redis.Nil
}

func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	return r.UniversalClient.LPush(key).Result()
}

func (r *Redis) LPushX(key string, value interface{}) (int64, error) {
	return r.UniversalClient.LPushX(key, value).Result()
}

func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return r.UniversalClient.LRange(key, start, stop).Result()
}

func (r *Redis) LRem(key string, count int64, value interface{}) (int64, error) {
	return r.UniversalClient.LRem(key, count, value).Result()
}

func (r *Redis) LSet(key string, index int64, value interface{}) (string, error) {
	return r.UniversalClient.LSet(key, index, value).Result()
}

func (r *Redis) LTrim(key string, start, stop int64) (string, error) {
	return r.UniversalClient.LTrim(key, start, stop).Result()
}

func (r *Redis) RPop(key string) (string, error) {
	return r.UniversalClient.RPop(key).Result()
}

func (r *Redis) RPopLPush(source, destination string) (string, error) {
	return r.UniversalClient.RPopLPush(source, destination).Result()
}

func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {
	return r.UniversalClient.RPush(key, values...).Result()
}

func (r *Redis) RPushX(key string, value interface{}) (int64, error) {
	return r.UniversalClient.RPushX(key, value).Result()
}

func (r *Redis) SAdd(key string, members ...interface{}) (int64, error) {
	return r.UniversalClient.SAdd(key, members...).Result()
}

func (r *Redis) SCard(key string) (int64, error) {
	return r.UniversalClient.SCard(key).Result()
}

func (r *Redis) SDiff(keys ...string) ([]string, error) {
	return r.UniversalClient.SDiff(keys...).Result()
}

func (r *Redis) SDiffStore(destination string, keys ...string) (int64, error) {
	return r.UniversalClient.SDiffStore(destination, keys...).Result()
}

func (r *Redis) SInter(keys ...string) ([]string, error) {
	return r.UniversalClient.SInter(keys...).Result()
}

func (r *Redis) SInterStore(destination string, keys ...string) (int64, error) {
	return r.UniversalClient.SInterStore(destination, keys...).Result()
}

func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return r.UniversalClient.SIsMember(key, member).Result()
}

func (r *Redis) SMembers(key string) ([]string, error) {
	return r.UniversalClient.SMembers(key).Result()
}

func (r *Redis) SMembersMap(key string) (map[string]struct{}, error) {
	return r.UniversalClient.SMembersMap(key).Result()
}

func (r *Redis) SMove(source, destination string, member interface{}) (bool, error) {
	return r.UniversalClient.SMove(source, destination, member).Result()
}

func (r *Redis) SPop(key string) (string, error) {
	val, err := r.UniversalClient.SPop(key).Result()
	if err != redis.Nil {
		return val, err
	}
	return val, redis.Nil
}

func (r *Redis) SPopN(key string, count int64) ([]string, error) {
	return r.UniversalClient.SPopN(key, count).Result()
}

func (r *Redis) SRandMember(key string) (string, error) {
	val, err := r.UniversalClient.SRandMember(key).Result()
	if err != redis.Nil {
		return val, err
	}
	return val, redis.Nil
}

func (r *Redis) SRandMemberN(key string, count int64) ([]string, error) {
	return r.UniversalClient.SRandMemberN(key, count).Result()
}

func (r *Redis) SRem(key string, members ...interface{}) (int64, error) {
	return r.UniversalClient.SRem(key, members...).Result()
}

func (r *Redis) SUnion(keys ...string) ([]string, error) {
	return r.UniversalClient.SUnion(keys...).Result()
}

func (r *Redis) SUnionStore(destination string, keys ...string) (int64, error) {
	return r.UniversalClient.SUnionStore(destination, keys...).Result()
}

func (r *Redis) ZCard(key string) (int64, error) {
	return r.UniversalClient.ZCard(key).Result()
}

func (r *Redis) ZAdd(key string, member ...Z) (int64, error) {
	return r.UniversalClient.ZAdd(key, member...).Result()
}

func (r *Redis) ZCount(key, min, max string) (int64, error) {
	return r.UniversalClient.ZCount(key, min, max).Result()
}

func (r *Redis) ZLexCount(key, min, max string) (int64, error) {
	return r.UniversalClient.ZLexCount(key, min, max).Result()
}

func (r *Redis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return r.UniversalClient.ZIncrBy(key, increment, member).Result()
}

func (r *Redis) ZInterStore(destination string, store ZStore, keys ...string) (int64, error) {
	return r.UniversalClient.ZInterStore(destination, store, keys...).Result()
}

func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	return r.UniversalClient.ZRange(key, start, stop).Result()
}

func (r *Redis) ZRangeByScore(key, min, max string, offset, count int64) ([]string, error) {
	return r.UniversalClient.ZRangeByScore(key, redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (r *Redis) ZRangeByLex(key, min, max string, offset, count int64) ([]string, error) {
	return r.UniversalClient.ZRangeByLex(key, redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (r *Redis) ZRank(key, member string) (int64, error) {
	return r.UniversalClient.ZRank(key, member).Result()
}

func (r *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	return r.UniversalClient.ZRem(key, members...).Result()
}

func (r *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return r.UniversalClient.ZRemRangeByRank(key, start, stop).Result()
}

func (r *Redis) ZRemRangeByScore(key, min, max string) (int64, error) {
	return r.UniversalClient.ZRemRangeByScore(key, min, max).Result()
}

func (r *Redis) ZRemRangeByLex(key, min, max string) (int64, error) {
	return r.UniversalClient.ZRemRangeByLex(key, min, max).Result()
}

func (r *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	return r.UniversalClient.ZRevRange(key, start, stop).Result()
}

func (r *Redis) ZRevRangeByScore(key, min, max string, offset, count int64) ([]string, error) {
	return r.UniversalClient.ZRevRangeByScore(key, redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (r *Redis) ZRevRangeByLex(key, min, max string, offset, count int64) ([]string, error) {
	return r.UniversalClient.ZRevRangeByLex(key, redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (r *Redis) ZRevRank(key, member string) (int64, error) {
	return r.UniversalClient.ZRevRank(key, member).Result()
}

func (r *Redis) ZScore(key, member string) (float64, error) {
	return r.UniversalClient.ZScore(key, member).Result()
}

func (r *Redis) ZUnionStore(dest string, weights []float64, aggregate string, keys ...string) (int64, error) {
	return r.UniversalClient.ZUnionStore(dest, redis.ZStore{
		Weights:   weights,
		Aggregate: aggregate,
	}, keys...).Result()
}

func (r *Redis) PFAdd(key string, els ...interface{}) (int64, error) {
	return r.UniversalClient.PFAdd(key, els...).Result()
}

func (r *Redis) PFCount(keys ...string) (int64, error) {
	return r.UniversalClient.PFCount(keys...).Result()
}

func (r *Redis) PFMerge(dest string, keys ...string) (string, error) {
	return r.UniversalClient.PFMerge(dest, keys...).Result()
}

func (r *Redis) HGObject(key string, out interface{}) error {
	data, err := r.UniversalClient.HGetAll(key).Result()
	if err != redis.Nil {
		if err == redis.Nil {
			return err
		}
		return errors.WithStack(err)
	}
	return errors.WithStack(xmap.Unmarshal(data, out))
}

func (r *Redis) HSObject(key string, in interface{}) error {
	err := r.UniversalClient.HMSet(key, xmap.Marshal(in)).Err()
	if err != redis.Nil {
		return errors.WithStack(err)
	}
	return redis.Nil
}

func (r *Redis) HMGObject(key string, out interface{}, fields ...string) error {
	vals, err := r.UniversalClient.HMGet(key, fields...).Result()
	if err != redis.Nil {
		if err == redis.Nil {
			return err
		}
		return errors.WithStack(err)
	}
	data := make(map[string]interface{}, len(fields))
	for i, field := range fields {
		data[field] = vals[i]
	}
	return errors.WithStack(mapstructure.Decode(data, out))
}
