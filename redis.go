package shield

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type RedisStore struct {
	redis      redis.Conn
	addr       string
	password   string
	sumKey     string
	classKey   string
	classesKey string
}

func NewRedisStore(addr, password string) Store {
	return &RedisStore{
		addr:       addr,
		password:   password,
		sumKey:     "shield:sum",
		classKey:   "shield:class",
		classesKey: "shield:classes",
	}
}

func (rs *RedisStore) conn() (conn redis.Conn, err error) {
	if rs.redis == nil {
		c, err2 := redis.DialTimeout("tcp", rs.addr, 0, 1*time.Second, 1*time.Second)
		if err2 != nil {
			return nil, err2
		}
		if rs.password != "" {
			_, authErr := redis.String(c.Do("AUTH", rs.password))
			if authErr != nil {
				err = authErr
				return
			}
		}
		rs.redis = c
	}
	return rs.redis, nil
}

func (rs *RedisStore) Classes() (a []string, err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	classes, err := redis.Strings(c.Do("SMEMBERS", rs.classesKey))
	if err != nil {
		return
	}
	for _, v := range classes {
		a = append(a, string(v))
	}
	return
}

func (rs *RedisStore) AddClass(class string) (err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	if class == "" {
		panic("invalid class: " + class)
	}
	_, err = c.Do("SADD", rs.classesKey, []byte(class))
	return err
}

func (rs *RedisStore) ClassWordCount(class, word string) (i int64, err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	b, err := redis.Bytes(c.Do("HGET", rs.classKey+":"+class, word))
	if err != nil {
		return
	}
	i, err = strconv.ParseInt(string(b), 10, 64)
	return
}

func (rs *RedisStore) IncrementClassWordCounts(m map[string]map[string]int64) (err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	c.Send("MULTI")
	for class, words := range m {
		for word, count := range words {
			c.Send("HINCRBY", rs.classKey+":"+class, word, count)
			c.Send("HINCRBY", rs.sumKey, class, count)
		}
	}
	c.Send("EXEC")
	err = c.Flush()
	if err != nil {
		return
	}
	_, err = c.Receive()
	return
}

func (rs *RedisStore) TotalClassWordCounts() (m map[string]int64, err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	values, err := redis.Values(c.Do("HGETALL", rs.sumKey))
	if err != nil {
		return
	}
	m = make(map[string]int64)
	for len(values) > 0 {
		var class string
		var count int64
		values, err = redis.Scan(values, &class, &count)
		if err != nil {
			return
		}
		m[class] = count
	}
	return
}

func (rs *RedisStore) Reset() (err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	a, err := redis.Strings(c.Do("KEYS", "shield:*"))
	if err != nil {
		return
	}
	c.Send("MULTI")
	for _, key := range a {
		c.Send("DEL", key)
	}
	c.Send("EXEC")
	err = c.Flush()
	if err != nil {
		return
	}
	_, err = c.Receive()
	return
}
