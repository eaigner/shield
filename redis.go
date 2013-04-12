package shield

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisStore struct {
	redis      redis.Conn
	addr       string
	password   string
	sumKey     string
	classKey   string
	classesKey string
	logger     *log.Logger
	prefix     string
}

func NewRedisStore(addr, password string, logger *log.Logger, prefix string) Store {
	return &RedisStore{
		addr:       addr,
		password:   password,
		sumKey:     "shield:sum",
		classKey:   "shield:class",
		classesKey: "shield:classes",
		logger:     logger,
		prefix:     prefix,
	}
}

func (rs *RedisStore) conn() (conn redis.Conn, err error) {
	if rs.redis == nil {
		c, err2 := redis.Dial("tcp", rs.addr)
		if err2 != nil {
			return nil, err2
		}
		if rs.logger != nil {
			c = redis.NewLoggingConn(c, rs.logger, rs.prefix)
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
	_, err = c.Do("SADD", rs.classesKey, class)
	return err
}

func (rs *RedisStore) ClassWordCounts(class string, words []string) (mc map[string]int64, err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	c.Send("MULTI")
	key := rs.classKey + ":" + class
	args := make([]interface{}, 0, len(words)+1)
	args = append(args, key)
	for _, v := range words {
		args = append(args, v)
	}
	c.Send("HMGET", args...)
	values, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		return
	}

	if len(values) > 0 {
		if x, ok := values[0].([]interface{}); ok {
			values = x
		}
	}

	var i int64

	mc = make(map[string]int64)
	for len(values) > 0 {
		var count int64
		values, err = redis.Scan(values, &count)
		if err != nil {
			return
		}
		mc[words[i]] = count
		i++
	}
	return
}

func (rs *RedisStore) IncrementClassWordCounts(m map[string]map[string]int64) (err error) {
	c, err := rs.conn()
	if err != nil {
		return
	}
	type keyPath struct {
		class string
		word  string
	}
	var wasDec bool
	paths := make([]*keyPath, 0, len(m)*100)
	c.Send("MULTI")
	for class, words := range m {
		for word, d := range words {
			c.Send("HINCRBY", rs.classKey+":"+class, word, d)
			c.Send("HINCRBY", rs.sumKey, class, d)
			paths = append(paths, &keyPath{
				class: class,
				word:  word,
			})
			if d < 0 {
				wasDec = true
			}
		}
	}
	values, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		return
	}

	// If we decrement something, we have to check if we went
	// below 0 afterwards and reset to 0 if necessary
	if wasDec {
		c.Send("MULTI")
		for i := 0; i < len(values); i += 2 {
			if v := values[i].(int64); v < 0 {
				kp := paths[i/2]
				d := v * -1
				c.Send("HINCRBY", rs.classKey+":"+kp.class, kp.word, d)
				c.Send("HINCRBY", rs.sumKey, kp.class, d)
			}
		}
		_, err = c.Do("EXEC")
	}
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
	_, err = c.Do("EXEC")
	return
}
