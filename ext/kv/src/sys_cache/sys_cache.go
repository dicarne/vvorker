package sys_cache

import (
	"errors"
	"fmt"
	"vvorker/utils"

	"github.com/nutsdb/nutsdb"
	"github.com/sirupsen/logrus"
)

var db *nutsdb.DB

func InitCache(_db *nutsdb.DB) {
	db = _db
}

func Put(key string, value []byte, ttl int) (int, error) {
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Put("sys_cache", []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func Get(key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *nutsdb.Tx) error {
		v, err := tx.Get("sys_cache", []byte(key))
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		if errors.Is(err, nutsdb.ErrKeyNotFound) {
			return []byte(""), nil
		}
		return nil, err
	}
	return value, nil
}

func Del(key string) error {
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete("sys_cache", []byte(key))
	})
}

func PutNX(key string, value []byte, ttl int) (int, error) {
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.PutIfNotExists("sys_cache", []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func PutXX(key string, value []byte, ttl int) (int, error) {
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.PutIfExists("sys_cache", []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func GlobalCache(key string, getValueFn func() ([]byte, error), ttl int) ([]byte, error) {
	cacheKey := fmt.Sprintf("db:workerd:%s_cache:", key)
	lockKey := fmt.Sprintf("db:workerd:%s_lock:", key)

	if v, err := Get(cacheKey); err == nil && len(v) != 0 {
		go func() {
			guid := utils.GenerateUID()
			_, _ = PutNX(lockKey, []byte(guid), 10)
			guid2, _ := Get(lockKey)
			if guid == string(guid2) {
				v, err := getValueFn()
				if err != nil {
					Del(lockKey)
					Del(cacheKey)
					return
				}
				Put(cacheKey, v, 0)
			}
		}()
		return v, nil
	}

	v, err := getValueFn()
	if err != nil {
		return nil, err
	}
	Put(cacheKey, v, 0)
	Put(lockKey, []byte(utils.GenerateUID()), ttl)
	return v, nil
}
