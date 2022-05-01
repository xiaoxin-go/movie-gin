package utils

import (
	"encoding/json"
	"fmt"
	config "movie/conf"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	Network string
	Host    string
	Port    string
	conn    redis.Conn
	DB      int
	Error   error
}

func NewRedisDefault() Redis {
	return NewRedis("tcp", config.Config.Redis.Host, config.Config.Redis.Port, 0)
}

// 初始化redis, 连接方式、主机名加端口、db
func NewRedis(network, host, port string, db int) (r Redis) {
	r = Redis{Network: network, Host: host, Port: port, DB: db}
	r.Conn()
	return
}

// 连接redis
func (r *Redis) Conn() {
	r.conn, r.Error = redis.Dial(r.Network, fmt.Sprintf("%s:%s", r.Host, r.Port), redis.DialDatabase(r.DB))
}

// 设置字符串
func (r *Redis) SetString(key, value string) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	_, r.Error = r.conn.Do("SET", key, value)
}

// 获取字符串
func (r *Redis) GetString(key string) (result string) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	var bytes []byte
	bytes, r.Error = redis.Bytes(r.conn.Do("GET", key))
	result = string(bytes)
	return
}

// 设置Json类型
func (r *Redis) SetJson(key string, value interface{}) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	var bytes []byte
	bytes, r.Error = json.Marshal(value)
	if r.Error != nil {
		return
	}
	_, r.Error = r.conn.Do("SET", key, bytes)
}

// 获取json类型
func (r *Redis) GetJson(key string, value interface{}) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	var bytes []byte
	bytes, r.Error = redis.Bytes(r.conn.Do("GET", key))
	if r.Error != nil {
		return
	}
	r.Error = json.Unmarshal(bytes, value)
}

// 设置过期时间
func (r *Redis) SetExpire(key string, expire int64) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	_, r.Error = r.conn.Do("EXPIRE", key, expire)
}

// 删除key
func (r *Redis) Delete(key string) {
	if r.conn == nil {
		r.Conn()
	}
	if r.Error != nil {
		return
	}
	_, r.Error = r.conn.Do("DEL", key)
}

// 关闭
func (r *Redis) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
