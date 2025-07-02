package bolt

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// DB represents a BoltDB database wrapper
type DB struct {
	Path string
	db   *bolt.DB
}

// Open opens the BoltDB database
func (b *DB) Open() error {
	db, err := bolt.Open(b.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open db: %v", err)
	}
	b.db = db
	return nil
}

// Close closes the BoltDB database
func (b *DB) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

// GetBuckets returns all buckets in the database
func (b *DB) GetBuckets() ([]string, error) {
	var buckets []string
	err := b.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	return buckets, err
}

// GetKeysInBucket returns all keys in a bucket
func (b *DB) GetKeysInBucket(bucketName string) ([]string, error) {
	var keys []string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		return bucket.ForEach(func(k, _ []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})
	return keys, err
}

// GetValue returns the value for a key in a bucket
func (b *DB) GetValue(bucketName, key string) ([]byte, error) {
	var value []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		value = bucket.Get([]byte(key))
		return nil
	})
	return value, err
}

// CreateBucket creates a new bucket
func (b *DB) CreateBucket(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

// DeleteBucket deletes a bucket
func (b *DB) DeleteBucket(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

// PutValue puts a value for a key in a bucket
func (b *DB) PutValue(bucketName, key string, value []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		return bucket.Put([]byte(key), value)
	})
}

// DeleteValue deletes a key from a bucket
func (b *DB) DeleteValue(bucketName, key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		return bucket.Delete([]byte(key))
	})
}
