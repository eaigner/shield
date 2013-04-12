package shield

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
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
	type tuple struct {
		word string
		d    int64
	}
	decrTuples := make(map[string][]*tuple, len(m))

	// Apply only positive increments
	c.Send("MULTI")
	for class, words := range m {
		for word, d := range words {
			if d > 0 {
				c.Send("HINCRBY", rs.classKey+":"+class, word, d)
				c.Send("HINCRBY", rs.sumKey, class, d)
			} else {
				decrTuples[class] = append(decrTuples[class], &tuple{
					word: word,
					d:    d,
				})
			}
		}
	}
	_, err = redis.Values(c.Do("EXEC"))
	if err != nil {
		return
	}

	// If we decrement something, we have to check if we are
	// about to drop below 0 and adjust the value accordingly.
	//
	// TODO: This isn't terribly performant because we have to
	// to 2 trips per value, try to optimize some time.
	//
	for class, paths := range decrTuples {
		key := rs.classKey + ":" + class

		// Build HMGET params
		hmget := make([]interface{}, 0, len(paths))
		hmget = append(hmget, key)
		for _, path := range paths {
			hmget = append(hmget, path.word)
		}

		values, err2 := redis.Strings(c.Do("HMGET", hmget...))
		if err2 != nil {
			return
		}

		c.Send("MULTI")
		for i, v := range values {
			path := paths[i]
			x, err2 := strconv.ParseInt(v, 10, 64)
			if err2 != nil {
				panic(err2)
			}
			if x != 0 {
				d := path.d
				if (x + d) < 0 {
					d = x * -1
				}
				c.Send("HINCRBY", key, path.word, d)
				c.Send("HINCRBY", rs.sumKey, class, d)
			}
		}
		_, err2 = c.Do("EXEC")
		if err2 != nil {
			return err2
		}
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
