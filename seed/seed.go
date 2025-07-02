package main

import (
	"log"

	"github.com/boltdb/bolt"
)

func seed(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("bucket1"))
		if err != nil {
			return err
		}
		bucket := tx.Bucket([]byte("bucket1"))
		err = bucket.Put([]byte("key1"), []byte("value1"))
		if err != nil {
			return err
		}
		return nil
	})
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = seed(db)
	if err != nil {
		log.Fatal(err)
	}
}
