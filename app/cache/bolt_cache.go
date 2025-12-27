package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

// BoltCache BoltDB缓存结构
type BoltCache struct {
	db     *bbolt.DB
	dbPath string
	mutex  sync.Mutex // 简化锁机制，适应单用户场景
}

// CacheItem 缓存项
type CacheItem struct {
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"created_at"`
}

var (
    boltCache  *BoltCache
    bucketName = []byte("cache")
)

// InitAppCache 初始化应用缓存（可选项目目录）
// 若 projectDir 为空，则使用默认的 ./data 目录
func InitAppCache(projectDir string) error {
    // 获取当前工作目录
    workDir, err := os.Getwd()
    if err != nil {
        workDir = "."
    }

    // 优先使用传入的项目目录；为空则使用默认的 ./data 目录
    cacheDir := projectDir
    if cacheDir == "" {
        cacheDir = filepath.Join(workDir, "data")
    } else if !filepath.IsAbs(cacheDir) {
        cacheDir = filepath.Join(workDir, cacheDir)
    }

    // 目标DB路径
    targetDBPath := filepath.Join(cacheDir, "app_cache.bolt")

    // 如果已初始化且路径相同，直接返回
    if boltCache != nil && boltCache.db != nil && boltCache.dbPath == targetDBPath {
        return nil
    }

    // 若已存在但路径不同或未初始化：关闭旧连接，重新初始化
    if boltCache != nil && boltCache.db != nil {
        _ = boltCache.db.Close()
        boltCache = nil
    }

    // 确保cache目录存在
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return err
    }

    // 打开BoltDB数据库
    db, err := bbolt.Open(targetDBPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
    if err != nil {
        return err
    }

    // 创建bucket
    if err := db.Update(func(tx *bbolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(bucketName)
        return err
    }); err != nil {
        _ = db.Close()
        return err
    }

    boltCache = &BoltCache{db: db, dbPath: targetDBPath}
    return nil
}

// GetBoltCache 获取BoltDB缓存实例
func GetBoltCache() *BoltCache {
    // 不再在此自动初始化，避免生成默认路径的缓存文件
    return boltCache
}

// Close 关闭数据库连接
func Close() error {
	cache := GetBoltCache()
	if cache != nil && cache.db != nil {
		return cache.db.Close()
	}
	return nil
}

// Set 设置缓存值
func Set(key string, value interface{}) error {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return nil
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := CacheItem{
		Value:     value,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return cache.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		return bucket.Put([]byte(key), data)
	})
}

// Get 获取缓存值
func Get(key string, dest interface{}) bool {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return false
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var data []byte
	err := cache.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		data = bucket.Get([]byte(key))
		return nil
	})

	if err != nil || data == nil {
		return false
	}

	var item CacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return false
	}

	// 使用类型断言将值赋给目标变量
	switch d := dest.(type) {
	case *string:
		if v, ok := item.Value.(string); ok {
			*d = v
			return true
		}
	case *int:
		// JSON解析数字时通常为float64，需要转换
		if v, ok := item.Value.(float64); ok {
			*d = int(v)
			return true
		}
	case *int64:
		if v, ok := item.Value.(float64); ok {
			*d = int64(v)
			return true
		}
	case *bool:
		if v, ok := item.Value.(bool); ok {
			*d = v
			return true
		}
	case *interface{}:
		*d = item.Value
		return true
	}

	return false
}

// Delete 删除缓存值
func Delete(key string) error {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return nil
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	return cache.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		return bucket.Delete([]byte(key))
	})
}

// Clear 清空所有缓存
func Clear() error {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return nil
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	return cache.db.Update(func(tx *bbolt.Tx) error {
		// 删除bucket并重新创建
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(bucketName)
		return err
	})
}

// Exists 检查键是否存在
func Exists(key string) bool {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return false
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var exists bool
	cache.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		data := bucket.Get([]byte(key))
		exists = data != nil
		return nil
	})

	return exists
}

// Keys 获取所有键
func Keys() []string {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return []string{}
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var keys []string
	cache.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		cursor := bucket.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	return keys
}

// Size 获取缓存大小
func Size() int {
	cache := GetBoltCache()
	if cache == nil || cache.db == nil {
		return 0
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var count int
	cache.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		stats := bucket.Stats()
		count = stats.KeyN
		return nil
	})

	return count
}
