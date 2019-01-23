package inventory

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var (
	boltConfig = &bolt.Options{Timeout: 1 * time.Second}
	bucket     = []byte("inventory")
	itemsKey   = []byte("items")
)

type boltBucket struct {
	filename string
}

func newBoltBucket(filename string) (*boltBucket, error) {
	bucket := &boltBucket{filename: filename}
	err := bucket.ensureBucket()
	return bucket, err
}

func (i *boltBucket) ensureBucket() error {
	db, err := i.openBoltDB()

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return err
}

func (i *boltBucket) get() ([]byte, error) {
	db, err := i.openBoltDB()

	if err != nil {
		return nil, nil
	}

	defer db.Close()

	var bucketData []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		data := b.Get(itemsKey)

		bucketData = make([]byte, len(data))
		copy(bucketData, data)

		return nil
	})

	return bucketData, err
}

func (i *boltBucket) put(data []byte) error {
	db, err := i.openBoltDB()

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(itemsKey, data)

		if err != nil {
			log.Printf("problem adding data to bucket %v", err)
			return nil
		}

		return nil
	})

	return err
}

func (i *boltBucket) openBoltDB() (*bolt.DB, error) {
	db, err := bolt.Open(i.filename, 0600, boltConfig)

	if err != nil {
		return nil, fmt.Errorf("problem opening db '%s', %+v", i.filename, err)
	}

	return db, nil
}
