package clRedis

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

type RedisObject struct {
	myredis   *redis.Client
	prefix    string
	noPrefix  bool
	isCluster bool
	isOrigin  bool // 原生输出结果
}

var RedisPool map[string]*RedisObject
var Locker sync.RWMutex

func init() {
	RedisPool = make(map[string]*RedisObject)
}

func New(addr, _password, _webSite string) (*RedisObject, error) {

	Locker.Lock()
	defer Locker.Unlock()

	val, find := RedisPool[addr+_webSite]
	if find {
		redisPing := val.myredis.Ping()
		if redisPing.Err() == nil {
			return val, nil
		}
		delete(RedisPool, addr+_webSite)
	}

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    _password,
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	clrd := &RedisObject{
		myredis:   client,
		prefix:    _webSite,
		isCluster: false,
	}

	RedisPool[addr+_webSite] = clrd
	return clrd, nil
}

// 连线关闭
func (this *RedisObject) Close() {

	if this.myredis != nil {
		_ = this.myredis.Close()
	}
}

// 测试连线
func (this *RedisObject) Ping() bool {
	redisPing := this.myredis.Ping()
	if redisPing.Err() == nil {
		return true
	}
	return false
}

// 删除
func (this *RedisObject) Del(key string) error {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	i := this.myredis.Del(keys)
	return i.Err()
}

// 删除
func (this *RedisObject) DelNoPrefix(key string) error {
	i := this.myredis.Del(key)
	return i.Err()
}

// 设置
func (this *RedisObject) Set(key string, val interface{}, expire int32) error {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	err := this.myredis.Set(keys, buildRedisValue(keys, uint32(expire), val),
		time.Duration(time.Second*time.Duration(expire))).Err()
	return err
}

// 获取指定的值
func (this *RedisObject) Get(key string) string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	resp := this.myredis.Get(keys)
	if this.isOrigin {
		return resp.Val()
	}
	result := checkRedisValid(keys, resp)
	if result == "" {
		this.myredis.Del(keys)
	}
	return result
}

// 获取最原始的数据
func (this *RedisObject) GetRaw(key string) string {
	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	resp := this.myredis.Get(keys)
	if resp == nil {
		return ""
	}
	return resp.Val()
}

// 获取指定的值 多語言用 允許空直不刪除
func (this *RedisObject) GetLang(key string) string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	resp := this.myredis.Get(keys)
	result := checkRedisValid(keys, resp)
	return result
}

func (this *RedisObject) GetNoPrefix(key string) string {

	keys := key
	resp := this.myredis.Get(keys)
	result := checkRedisValid(keys, resp)
	if result == "" {
		this.myredis.Del(keys)
	}
	return result
}

// 获取指定的json结构
func (this *RedisObject) GetJson(key string) *clJson.JsonStream {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	obj := this.myredis.Get(keys)
	return clJson.New([]byte(checkRedisValid(keys, obj)))
}

// 设置hash结构
func (this *RedisObject) HSet(key string, field string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	value = buildRedisValue(keys+field, expire, value)
	rest := this.myredis.HSet(keys, field, value)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		return false
	}

	return true
}

// 设置hash结构
func (this *RedisObject) SetEx(key string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	value = buildRedisValue(keys, expire, value)
	rest := this.myredis.Set(keys, value, time.Duration(expire)*time.Second)
	if rest == nil {
		return false
	}
	return rest.Val() == "OK"
}

// 设置hash结构的值(保存为json)
func (this *RedisObject) HSetJson(key string, field string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	value = buildRedisValue(keys+field, expire, value)
	rest := this.myredis.HSet(keys, field, value)

	if rest == nil {
		return false
	}

	return rest.Val()
}

// 获取hash结构
func (this *RedisObject) HGet(key string, field string) string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	resp := this.myredis.HGet(keys, field)
	if this.isOrigin {
		return resp.Val()
	}
	result := checkRedisValid(keys+field, resp)
	if result == "" {
		this.myredis.HDel(keys, field)
	}
	return result
}

func (this *RedisObject) Origin() *RedisObject {
	var _redis = RedisObject{}
	_redis = *this
	_redis.isOrigin = true
	return &_redis
}

func (this *RedisObject) HDel(key string, field string) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	resp := this.myredis.HDel(keys, field)
	return resp.Val() > 0
}

// 获取hash结构的值
func (this *RedisObject) HGetJson(key string, field string) *clJson.JsonStream {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	val := this.myredis.HGet(keys, field)
	res := checkRedisValid(keys+field, val)
	if res == "" {
		return nil
	}
	return clJson.New([]byte(res))
}

// 获取全部的key
func (this *RedisObject) HGetKeys(key string, prefix string) []string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	val := this.myredis.HKeys(keys)
	if val == nil {
		return []string{}
	}
	resp := make([]string, 0)
	if prefix != "" {
		for _, row := range val.Val() {
			if strings.HasPrefix(row, prefix) {
				resp = append(resp, row)
			}
		}
	} else {
		resp = val.Val()
	}
	return resp
}

// 删除指定开头的keys
func (this *RedisObject) HDelKeys(key string, prefix string) {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	keylist := this.HGetKeys(keys, prefix)
	if len(keylist) > 0 {
		this.myredis.HDel(keys, keylist...)
	}
}

// 获取全部的hash字段
func (this *RedisObject) HGetAll(key string) map[string]string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	val := this.myredis.HGetAll(keys)
	return checkRedisValidMap(keys, val)
}

// 设置锁
func (this *RedisObject) SetNx(key string, value interface{}, expire uint32) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	value = buildRedisValue(keys, expire, value)
	rest := this.myredis.SetNX(keys, value, time.Duration(expire)*time.Second)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		fmt.Printf(">> SetNX |%v| Failed! Err:%v\n", keys, err)
		return false
	}

	return rest.Val()
}

// md5加密
func md5Encode(_data []byte) string {
	h := md5.New()
	h.Write(_data) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkRedisValidMap(keys string, targetData *redis.StringStringMapCmd) map[string]string {
	if targetData == nil || len(targetData.Val()) == 0 {
		return nil
	}

	resp := make(map[string]string)

	for key, val := range targetData.Val() {
		js := clJson.New([]byte(val))
		expireTime := js.GetUint32("expire")
		// 缓存到期
		if expireTime > 0 && expireTime < uint32(time.Now().Unix()) {
			continue
		}
		sign := md5Encode([]byte("Cache:__" + keys + key))
		if js.GetStr("sign") != sign {
			continue
		}
		resp[key] = js.GetStr("data")
	}
	return resp
}

// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkStringValid(keys string, _targetStr string) string {

	js := clJson.New([]byte(_targetStr))
	if js == nil {
		return ""
	}

	expireTime := js.GetUint32("expire")
	addtime := js.GetUint32("addtime")

	// 缓存到期  新增添加时间大于当前时间表示有问题
	if expireTime == 0 || expireTime < uint32(time.Now().Unix()) {
		return ""
	}
	if addtime > uint32(time.Now().Unix()) {
		return ""
	}

	sign := md5Encode([]byte("Cache:__" + keys))
	if js.GetStr("sign") != sign {
		return ""
	}

	return js.GetStr("data")
}

// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkRedisValid(keys string, targetData *redis.StringCmd) string {
	if targetData == nil || targetData.Val() == "" {
		return ""
	}

	return checkStringValid(keys, targetData.Val())
}

// 组装缓存的值
func buildRedisValue(keys string, expire uint32, data interface{}) string {
	return clJson.CreateBy(clJson.M{
		"data":    data,
		"addtime": uint32(time.Now().Unix()), // 写入redis缓存的时间
		"expire":  uint32(time.Now().Unix()) + expire,
		"sign":    md5Encode([]byte("Cache:__" + keys)),
	}).ToStr()
}

// 操作list结构 lpush,push 是会更新key的过期时间
func (this *RedisObject) Lpush(key string, expire uint32, values ...interface{}) bool {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}
	for k, value := range values {
		values[k] = buildRedisValue(keys, uint32(expire), value)
	}
	rest := this.myredis.LPush(keys, values...)

	// 设置过期时间
	if expire > 0 {
		this.myredis.Expire(keys, time.Duration(expire)*1000*time.Millisecond)
	}

	if rest == nil {
		return false
	}
	return true
}

// 操作list结构 lpop
func (this *RedisObject) Lpop(key string) string {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	val := this.myredis.LPop(keys)
	result := checkRedisValid(keys, val)
	return result
}

// 操作list结构 blpop
// 要操作的key
// 要等待的时间
func (this *RedisObject) LPOPWait(key string, _timeOut uint32) (error, []string) {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	val := this.myredis.BLPop(time.Duration(_timeOut)*time.Second, keys)
	if val.Err() != nil {
		if val.Err().Error() == "redis: nil" {
			return nil, nil
		}
		clLog.Debug("valErr: %v", val.Err().Error())
		return val.Err(), nil
	}

	res, err := val.Result()
	if err != nil {
		clLog.Debug("err: %v", err)
		return err, nil
	}

	result := make([]string, 0)
	for _, str := range res {
		isOkStr := checkStringValid(keys, str)
		if isOkStr != "" {
			result = append(result, isOkStr)
		}
	}
	return nil, result
}

// 取队列元素个数
func (this *RedisObject) Llen(key string) (error, int64) {
	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	result := this.myredis.LLen(keys)
	return result.Err(), result.Val()
}

// 操作list结构 rpop
func (this *RedisObject) Rpop(key string) interface{} {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	val := this.myredis.RPop(keys)
	result := checkRedisValid(keys, val)
	return result
}

// 删除list
func (this *RedisObject) DelList(key string) {

	keys := key
	if this.prefix != "" && !this.noPrefix {
		keys = this.prefix + "_" + key
	}

	this.myredis.LTrim(keys, 1, 0)
}

// 获取key列表
func (this *RedisObject) GetKeys(key string) []string {
	res := this.myredis.Keys(key)
	return res.Val()
}

// 删除所有的类似的key
func (this *RedisObject) DelAll(key string) {

	res := this.myredis.Keys(key)

	klist, _ := res.Result()
	this.myredis.Del(klist...)
}

// 判断键是否存在
func (this *RedisObject) IsExists(key string) bool {
	res := this.myredis.Exists(key)
	return res.Val() == 1
}

// 添加一个值
func (this *RedisObject) SetNXInt(key string, _val int64) bool {
	var res *redis.BoolCmd
	res = this.myredis.SetNX(key, _val, 0)
	return res.Val()
}

// 添加一个值
func (this *RedisObject) Increment(key string, _val int64) int64 {
	var res *redis.IntCmd
	if _val < 0 {
		res = this.myredis.DecrBy(key, -_val)
	} else {
		res = this.myredis.IncrBy(key, _val)
	}
	return res.Val()
}

func (this *RedisObject) HInc(key, field string, incr int64) int64 {
	cmd := this.myredis.HIncrBy(key, field, incr)
	return cmd.Val()
}

// 添加一个值
func (this *RedisObject) SetExpire(_key string, _second uint32) bool {
	var res *redis.BoolCmd

	res = this.myredis.Expire(_key, time.Second*time.Duration(_second))
	return res.Val()
}

func (this *RedisObject) Pointer() *redis.Client {
	return this.myredis
}

func (this *RedisObject) Keys(pattern string) []string {
	cmd := this.myredis.Keys(pattern)
	if cmd == nil {
		return nil
	}
	res, err := cmd.Result()
	if err != nil {
		clLog.Error("error:%v", err)
	}
	return res
}

func (this *RedisObject) NoPrefix(noprefix bool) *RedisObject {
	this.noPrefix = noprefix
	return this
}
