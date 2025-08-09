package kvnutsdb

import (
	"errors"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/ext/kv/src/sys_cache"

	"github.com/nutsdb/nutsdb"
	"github.com/sirupsen/logrus"
)

var db *nutsdb.DB

var buckets *defs.SyncMap[string, bool]

type KVNutsDB struct {
}

func init() {
	db2, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(conf.AppConfigInstance.LocalKVDir), // 数据库会自动创建这个目录文件
	)
	db = db2
	if err != nil {
		logrus.Panic(err)
	}
	buckets = defs.NewSyncMap(map[string]bool{})
	sys_cache.InitCache(db)
}

func (r *KVNutsDB) Close() {
	if db != nil {
		db.Close()
	}
}

func ExistBucket(bucket string) error {
	if _, exist := buckets.Get(bucket); exist {
		return nil
	}
	return db.Update(func(tx *nutsdb.Tx) error {
		tx.NewKVBucket(bucket)
		buckets.Set(bucket, true)
		return nil
	})
}

func (r *KVNutsDB) Put(bucket string, key string, value []byte, ttl int) (int, error) {
	ExistBucket(bucket)
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Put(bucket, []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func (r *KVNutsDB) PutNX(bucket string, key string, value []byte, ttl int) (int, error) {
	ExistBucket(bucket)
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.PutIfNotExists(bucket, []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func (r *KVNutsDB) PutXX(bucket string, key string, value []byte, ttl int) (int, error) {
	ExistBucket(bucket)
	code := 0
	return code, db.Update(func(tx *nutsdb.Tx) error {
		err := tx.PutIfExists(bucket, []byte(key), value, uint32(ttl))
		if err != nil {
			code = 1
			logrus.Error(err)
		}
		return err
	})
}

func (r *KVNutsDB) Get(bucket string, key string) ([]byte, error) {
	ExistBucket(bucket)
	var value []byte
	err := db.View(func(tx *nutsdb.Tx) error {
		v, err := tx.Get(bucket, []byte(key))
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		if errors.Is(err, nutsdb.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

func (r *KVNutsDB) Del(bucket string, key string) error {
	ExistBucket(bucket)
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	})
}

func (r *KVNutsDB) Keys(bucket string, prefix string, offset int, size int) ([]string, error) {
	ExistBucket(bucket)
	var result []string
	err := db.View(
		func(tx *nutsdb.Tx) error {
			prefixBytes := []byte(prefix)
			// Based on compiler feedback, tx.PrefixScan appears to return a slice of keys ([]byte),
			// not a slice of Entry structs.
			entries, err := tx.PrefixScan(bucket, prefixBytes, offset, size)
			if err != nil {
				return err
			}
			result = make([]string, 0, len(entries))
			for _, entry := range entries {
				result = append(result, string(entry))
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return result, nil
}
