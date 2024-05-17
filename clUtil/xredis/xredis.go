package xredis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clConfig"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type Option struct {
	Prefix   string
	Addr     string
	Password string
	Port     uint32
}
type XRedis struct {
	Option
	err   error // 错误内容
	redis *redis.Client
	ctx   context.Context
}

type SetValue struct {
	Value  interface{} `json:"v"`
	Expire int64       `json:"e"`
}

var client *XRedis

func NewWithOption(option Option) *XRedis {
	return newClient(option)
}

func New() *XRedis {
	option := Option{}
	option.Addr = clConfig.GetStr("REDIS_HOST", "")
	option.Prefix = clConfig.GetStr("REDIS_PREFIX", "")
	option.Password = clConfig.GetStr("REDIS_PASSWORD", "")
	option.Port = clConfig.GetUint32("REDIS_PORT", 6379)
	return NewWithOption(option)
}

func newClient(option Option) *XRedis {
	if client != nil {
		err := client.redis.Ping().Err()
		if err == nil {
			return client
		}
	}
	xredis := &XRedis{
		Option: option,
		redis:  nil,
		ctx:    nil,
	}
	var addr = xredis.Addr
	if xredis.Port > 0 {
		if !strings.Contains(addr, ":") {
			addr = fmt.Sprintf("%v:%v", addr, xredis.Port)
		}
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    xredis.Password,
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
	})
	if err := rdb.Ping().Err(); err != nil {
		clLog.Error("redis connect error:%v", err)
		return nil
	}
	xredis.redis = rdb
	client = xredis
	return xredis
}

// 右出
func (this *XRedis) RPop(key string) *StringValue {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return defaultStringValue
	}
	cmd := this.redis.RPop(this.buildKey(key))
	result := &StringValue{}
	res, err := cmd.Result()
	result.Value = res
	this.err = err
	return result
}

// 同步锁
func (this *XRedis) HSetNx(key, field string, value interface{}, expire int64) bool {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return false
	}
	cmd := this.redis.HExists(this.buildKey(key), field)
	if cmd.Val() {
		// 存在
		if !this.HGet(key, field).IsExpire() {
			// 未过期
			return false
		}
	}
	cmd = this.redis.HSetNX(this.buildKey(key), field, this.buildSetValue(value, expire))
	result, err := cmd.Result()
	this.err = err
	return result
}

// 同步锁
func (this *XRedis) SetNx(key string, value interface{}, expire int64) bool {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return false
	}
	cmd := this.redis.SetNX(this.buildKey(key), value, time.Second*time.Duration(expire))
	result, err := cmd.Result()
	this.err = err
	return result
}

// 左进
func (this *XRedis) LPush(key string, value ...interface{}) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.LPush(this.buildKey(key), value...)
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) HDecrBy(key, field string, value int64) error {
	return this.HIncrBy(this.buildKey(key), field, -value)
}
func (this *XRedis) HDecr(key, field string) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.HIncrBy(this.buildKey(key), field, -1)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) DecrBy(key string, value int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.DecrBy(this.buildKey(key), value)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) Decr(key string) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.Decr(this.buildKey(key))
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) HIncrByFloat(key, field string, value float64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.HIncrByFloat(this.buildKey(key), field, value)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) HIncrBy(key, field string, value int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.HIncrBy(this.buildKey(key), field, value)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) HIncr(key, field string) error {
	return this.HIncrBy(this.buildKey(key), field, 1)
}
func (this *XRedis) IncrByFloat(key string, value float64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.IncrByFloat(this.buildKey(key), value)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) IncrBy(key string, value int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.IncrBy(this.buildKey(key), value)
	this.err = cmd.Err()
	return this.err
}
func (this *XRedis) Incr(key string) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return this.err
	}
	cmd := this.redis.Incr(this.buildKey(key))
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) HGet(key, field string) *SetValue {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return defaultSetValue
	}
	cmd := this.redis.HGet(this.buildKey(key), field)
	result, err := cmd.Bytes()
	if err != nil {
		this.err = err
		return defaultSetValue
	}
	setValue := this.parseSetValue(result)
	if setValue.IsExpire() {
		this.HDel(key, field)
	}
	return setValue
}
func (this *XRedis) parseSetValue(value []byte) *SetValue {
	obj := defaultSetValue
	json.Unmarshal(value, obj)
	return obj
}

func (this *XRedis) HDel(key string, field ...string) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	cmd := this.redis.HDel(this.buildKey(key), field...)
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) Del(keys ...string) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	cmd := this.redis.Del(this.buildKeys(keys...)...)
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) HSet(key, field string, value interface{}, expire int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	cmd := this.redis.HSet(this.buildKey(key), field, this.buildSetValue(value, expire))
	this.err = cmd.Err()
	return this.err
}

func (this *XRedis) buildSetValue(value interface{}, expire int64) []byte {
	SetValue := SetValue{
		Value:  value,
		Expire: time.Now().Unix() + expire,
	}
	result, _ := json.Marshal(SetValue)
	return result
}

// 读取
func (this *XRedis) GetObj(key string, value interface{}) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Bytes()
	if err != nil {
		this.err = err
		return err
	}
	err = json.Unmarshal(result, value)
	this.err = err
	return err
}

// 读取错误信息
func (this *XRedis) Error() error {
	return this.err
}

func (this *XRedis) Float64(key string) float64 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Float64()
	if err != nil {
		this.err = err
	}
	return result
}
func (this *XRedis) Float32(key string) float32 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Float32()
	if err != nil {
		this.err = err
	}
	return result
}
func (this *XRedis) Uint(key string) uint {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Uint64()
	if err != nil {
		this.err = err
	}
	return uint(result)
}
func (this *XRedis) Int(key string) int {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Int()
	if err != nil {
		this.err = err
	}
	return result
}
func (this *XRedis) Int32(key string) int32 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Int64()
	if err != nil {
		this.err = err
	}
	return int32(result)
}
func (this *XRedis) Int64(key string) int64 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Int64()
	if err != nil {
		this.err = err
	}
	return result
}

func (this *XRedis) Uint64(key string) uint64 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Uint64()
	if err != nil {
		this.err = err
	}
	return result
}

func (this *XRedis) Client() *redis.Client {
	return this.redis
}

func (this *XRedis) GetCmd(key string) *redis.StringCmd {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return nil
	}
	cmd := this.redis.Get(this.buildKey(key))
	return cmd
}
func (this *XRedis) Uint32(key string) uint32 {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return 0
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Uint64()
	if err != nil {
		this.err = err
	}
	return uint32(result)
}

// 读取
func (this *XRedis) Get(key string) string {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return ""
	}
	cmd := this.redis.Get(this.buildKey(key))
	result, err := cmd.Result()
	if err != nil {
		this.err = err
	}
	return result
}

func (this *XRedis) Set(key string, value interface{}, expire int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	cmd := this.redis.Set(this.buildKey(key), value, time.Second*time.Duration(expire))
	this.err = cmd.Err()
	return cmd.Err()
}

func (this *XRedis) SetObj(key string, obj interface{}, expire int64) error {
	if this.redis == nil {
		this.err = fmt.Errorf("redis connect error")
		return fmt.Errorf("redis connect error")
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		this.err = err
		return err
	}
	cmd := this.redis.Set(this.buildKey(key), bytes, time.Second*time.Duration(expire))
	this.err = cmd.Err()
	return cmd.Err()
}

func (this *XRedis) buildKey(key string) string {
	if this.Prefix != "" {
		return fmt.Sprintf("%v_%v", this.Prefix, key)
	}
	return key
}

func (this *XRedis) buildKeys(_keys ...string) []string {
	keys := []string{}
	for _, key := range _keys {
		keys = append(keys, this.buildKey(key))
	}
	return keys
}
