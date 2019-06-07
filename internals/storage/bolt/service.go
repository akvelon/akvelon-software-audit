package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	DBPath     = "akvelonaudit.db"
	RepoBucket = "Repository"
)

type RepoScanResult struct {
	File       string `json:"File"`
	License    string `json:"License"`
	Confidence string `json:"Confidence"`
	Size       string `json:"Size"`
}

type NotFoundError struct {
	repo string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%q not found in cache", n.repo)
}

// InitStorage initialize DB for usage.
func InitStorage() error {
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RepoBucket))
		if err != nil {
			return err
		}
		return err
	})
	return err
}

// GetRepoFromDB returns repo data if exists.
func GetRepoFromDB(repo string) ([]RepoScanResult, error) {
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt database during GET: %v", err)
	}
	defer db.Close()
	var result []byte
	resp := []RepoScanResult{}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return errors.New("No repo bucket")
		}
		result = b.Get([]byte(repo))

		if result == nil {
			return NotFoundError{repo}
		}

		err = json.Unmarshal(result, &resp)
		if err != nil {
			return fmt.Errorf("failed to parse JSON for %q in result", repo)
		}

		return nil
	})

	if err != nil {
		switch err.(type) {
		case NotFoundError:
			// do nothing at that time
		default:
			log.Println("ERROR:", err) // log error, but continue
		}
		return nil, err
	}

	return resp, nil
}

// SaveRepoToDB save repo data.
func SaveRepoToDB(key string, data []byte) error {
	db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		log.Printf("Saving %q to db...", key)

		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("repo bucket not found")
		}

		// save repo to cache
		err = b.Put([]byte(key), data)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Bolt writing error:  %v", err)
	}
	return nil
}
