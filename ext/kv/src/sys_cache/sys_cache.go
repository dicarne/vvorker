package sys_cache

import (
	"errors"

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
