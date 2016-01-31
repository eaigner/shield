package shield

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	bolt       *bolt.DB
	path       string
	sumKey     string
	classKey   string
	classesKey string
	logger     *log.Logger
	prefix     string
}

func NewBoltStore(path string, logger *log.Logger, prefix string) *BoltStore {
	bs := &BoltStore{
		path:       path,
		sumKey:     "shield:sum",
		classKey:   "shield:class",
		classesKey: "shield:classes",
		logger:     logger,
		prefix:     prefix,
	}

	bs.init()
	return bs
}

type Bucket struct {
	*bolt.Bucket
}

func (b Bucket) Get(key string) int64 {
	buff := b.Bucket.Get([]byte(key))
	if val, n := binary.Varint(buff); n > 0 {
		return val
	}
	return 0
}

func (b Bucket) IncrementBy(key string, inc int64) error {
	ret := b.Bucket.Get([]byte(key))

	value := int64(0)
	if val, n := binary.Varint(ret); n > 0 {
		value = val
	}

	value += int64(inc)

	buff := make([]byte, 8)
	binary.PutVarint(buff, int64(value))

	if err := b.Bucket.Put([]byte(key), buff); err != nil {
		return err
	}

	return nil
}

func (b Bucket) Update(key string, value int64) error {
	buff := make([]byte, 8)
	binary.PutVarint(buff, value)
	return b.Bucket.Put([]byte(key), buff)
}

func (rs *BoltStore) init() (conn *bolt.DB, err error) {
	if rs.bolt == nil {
		db, err := bolt.Open(rs.path, 0600, nil)
		if err != nil {
			return nil, err
		}

		rs.bolt = db

		tx, err := db.Begin(true)
		if err != nil {
			return nil, err
		}

		defer tx.Rollback()

		if _, err := tx.CreateBucketIfNotExists([]byte(rs.sumKey)); err != nil {
			return nil, err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(rs.classKey)); err != nil {
			return nil, err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(rs.classesKey)); err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

	}
	return rs.bolt, nil
}

func (rs *BoltStore) Classes() (a []string, err error) {
	err = rs.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(rs.classesKey))

		return b.ForEach(func(k, v []byte) error {
			a = append(a, string(v))
			return nil
		})
	})

	return
}

func (rs *BoltStore) AddClass(class string) (err error) {
	if class == "" {
		return fmt.Errorf("invalid class: %s", class)
	}

	err = rs.bolt.Update(func(tx *bolt.Tx) error {
		b := Bucket{tx.Bucket([]byte(rs.classesKey))}
		return b.Update(class, 0)
	})

	return
}

func (rs *BoltStore) ClassWordCounts(class string, words []string) (mc map[string]int64, err error) {
	key := fmt.Sprintf("%s:%s", rs.classKey, class)

	if err = rs.bolt.Update(func(tx *bolt.Tx) error {
		b := Bucket{tx.Bucket([]byte(key))}

		mc = make(map[string]int64)
		for _, v := range words {
			mc[v] = b.Get(v)
		}

		return nil
	}); err != nil {
		return
	}

	return
}

func (rs *BoltStore) IncrementClassWordCounts(m map[string]map[string]int64) (err error) {
	type tuple struct {
		word string
		d    int64
	}

	decrTuples := make(map[string][]*tuple, len(m))

	if err = rs.bolt.Update(func(tx *bolt.Tx) error {
		sb := Bucket{tx.Bucket([]byte(rs.sumKey))}

		for class, words := range m {
			for word, d := range words {
				if d > 0 {
					key := fmt.Sprintf("%s:%s", rs.classKey, class)

					if bucket, err := tx.CreateBucketIfNotExists([]byte(key)); err == nil {
						b := Bucket{bucket}
						b.IncrementBy(word, d)
					}

					sb.IncrementBy(class, d)
				} else {
					decrTuples[class] = append(decrTuples[class], &tuple{
						word: word,
						d:    d,
					})
				}
			}
		}

		for class, paths := range decrTuples {
			key := fmt.Sprintf("%s:%s", rs.classKey, class)

			b := Bucket{tx.Bucket([]byte(key))}

			for _, path := range paths {
				if x := b.Get(path.word); x != 0 {
					d := path.d
					if (x + d) < 0 {
						d = x * -1
					}

					b.IncrementBy(path.word, d)
					sb.IncrementBy(class, d)
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return
}

func (rs *BoltStore) TotalClassWordCounts() (m map[string]int64, err error) {
	m = make(map[string]int64)

	err = rs.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(rs.sumKey))

		cursor := b.Cursor()
		for k, val := cursor.First(); k != nil; k, val = cursor.Next() {
			value, _ := binary.Varint(val)
			m[string(k)] = int64(value)
		}

		return nil
	})

	return
}

func (rs *BoltStore) Reset() (err error) {

	if rs.bolt != nil {
		rs.bolt.Close()

		defer os.Remove(rs.path)

		rs.bolt = nil
	}

	rs.init()
	return
}
