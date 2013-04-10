package shield

import (
	"github.com/eaigner/goredis"
	"strconv"
)

type RedisStore struct {
	client     *goredis.Client
	sumKey     string
	classKey   string
	classesKey string
}

func NewRedisStore(addr, password string, db int) Store {
	return &RedisStore{
		client: &goredis.Client{
			Addr:     addr,
			Db:       db,
			Password: password,
		},
		sumKey:     "shield:sum",
		classKey:   "shield:class",
		classesKey: "shield:classes",
	}
}

func (rs *RedisStore) Classes() (a []string, err error) {
	classes, err := rs.client.Smembers(rs.classesKey)
	if err != nil {
		return
	}
	for _, v := range classes {
		a = append(a, string(v))
	}
	return
}

func (rs *RedisStore) AddClass(class string) error {
	if class == "" {
		panic("invalid class: " + class)
	}
	_, err := rs.client.Sadd(rs.classesKey, []byte(class))
	return err
}

func (rs *RedisStore) ClassWordCount(class, word string) (i int64, err error) {
	b, err := rs.client.Hget(rs.classKey+":"+class, word)
	if err != nil {
		return
	}
	i, err = strconv.ParseInt(string(b), 10, 64)
	return
}

func (rs *RedisStore) IncrementClassWordCount(class, word string, i int64) (err error) {
	tx, err := rs.client.Transaction()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Discard()
		} else {
			_, err = tx.Exec()
		}
	}()
	_, err = tx.Hincrby(rs.classKey+":"+class, word, i)
	if err != nil {
		return
	}
	_, err = tx.Hincrby(rs.sumKey, class, i)
	return
}

func (rs *RedisStore) TotalClassWordCounts() (m map[string]int64, err error) {
	v := make(map[string]int64)
	if err = rs.client.Hgetall(rs.sumKey, &v); err == nil {
		m = v
	}
	return
}
