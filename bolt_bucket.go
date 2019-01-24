package cookme

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var (
	boltConfig = &bolt.Options{Timeout: 1 * time.Second}
	itemsKey   = []byte("items")
)

// BoltBucket is a wrapper around bolt just to store one blob of stuff in a bucket
type BoltBucket struct {
	filename string
	bucket   []byte
}

// NewBoltBucket creates a new BoltBucket, ensuring the bolt bucket is made
func NewBoltBucket(filename string, bucket string) (*BoltBucket, error) {
	bb := &BoltBucket{filename: filename, bucket: []byte(bucket)}
	err := bb.ensureBucket()
	return bb, err
}

// Get tries to retrieve the data inside the bucket
func (i *BoltBucket) Get() ([]byte, error) {
	db, err := i.openBoltDB()

	if err != nil {
		return nil, nil
	}

	defer db.Close()

	var bucketData []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(i.bucket)
		data := b.Get(itemsKey)

		bucketData = make([]byte, len(data))
		copy(bucketData, data)

		return nil
	})

	return bucketData, err
}

// Put will replace the data in the bucket
func (i *BoltBucket) Put(data []byte) error {
	db, err := i.openBoltDB()

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(i.bucket)
		err := b.Put(itemsKey, data)

		if err != nil {
			log.Printf("problem adding data to bucket %v", err)
			return nil
		}

		return nil
	})

	return err
}

func (i *BoltBucket) ensureBucket() error {
	db, err := i.openBoltDB()

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(i.bucket)
		return err
	})

	return err
}

func (i *BoltBucket) openBoltDB() (*bolt.DB, error) {
	db, err := bolt.Open(i.filename, 0600, boltConfig)

	if err != nil {
		return nil, fmt.Errorf("problem opening db '%s', %+v", i.filename, err)
	}

	return db, nil
}
