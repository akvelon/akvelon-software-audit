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
	// DBPath is the relative (or absolute) path to the bolt database file
	DBPath = "akvelonaudit.db"
	// RepoBucket is the bucket in which repos will be saved in the bolt DB
	RepoBucket = "Repository"
	// MetaBucket is the bucket containing meta info about db
	MetaBucket = "Meta"

	maxNumOfRecentItems = 5
)

// Storage keeps data in BoltDB
type Storage struct {
	db *bolt.DB
}

// NotFoundError is custom error for repositry
type NotFoundError struct {
	repo string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%q not found in cache", n.repo)
}

// OpenBoltDb opens bolt DB
func (b *Storage) OpenBoltDb() {
	var err error
	b.db, err = bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
}

// InitStorage initialize DB for usage.
func (b *Storage) InitStorage() error {
	b.OpenBoltDb()
	defer b.db.Close()

	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RepoBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(MetaBucket))
		return err
	})
	return err
}

// GetRecentlyViewed get list of recent items.
func (b *Storage) GetRecentlyViewed() ([]string, error) {
	b.OpenBoltDb()
	defer b.db.Close()

	recent := &[]recentItem{}
	err := b.db.View(func(tx *bolt.Tx) error {
		rb := tx.Bucket([]byte(MetaBucket))
		if rb == nil {
			return fmt.Errorf("meta bucket not found")
		}
		by := rb.Get([]byte("recent"))
		if by == nil {
			_, err := json.Marshal([]recentItem{})
			if err != nil {
				return err
			}
		}
		json.Unmarshal(by, recent)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Bolt writing error:  %v", err)
	}

	var recentRepos = make([]string, len(*recent))
	var j = len(*recent) - 1
	for _, r := range *recent {
		recentRepos[j] = r.Repo
		j--
	}
	return recentRepos, nil
}

// UpdateRecentlyViewed updated list of recent items.
func (b *Storage) UpdateRecentlyViewed(repo string) error {
	b.OpenBoltDb()
	defer b.db.Close()

	b.db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(MetaBucket))
		if mb == nil {
			return errors.New("meta bucket not found")
		}
		b := mb.Get([]byte("recent"))
		if b == nil {
			b, _ = json.Marshal([]recentItem{})
		}
		recent := []recentItem{}
		json.Unmarshal(b, &recent)

		for i := range recent {
			if recent[i].Repo == repo {
				return nil
			}
		}

		recent = append(recent, recentItem{Repo: repo})
		if len(recent) > maxNumOfRecentItems {
			// trim recent if it's grown to over maxNumOfRecentItems
			recent = (recent)[1 : maxNumOfRecentItems+1]
		}
		b, err := json.Marshal(&recent)
		if err != nil {
			return err
		}
		return mb.Put([]byte("recent"), b)
	})
	return nil
}

// GetRepoFromDB returns repo data if exists.
func (b *Storage) GetRepoFromDB(repo string) ([]byte, error) {
	b.OpenBoltDb()
	defer b.db.Close()
	var result []byte
	resp := []RepoScanResult{}

	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return errors.New("No repo bucket")
		}
		result = b.Get([]byte(repo))

		if result == nil {
			return NotFoundError{repo}
		}
		err := json.Unmarshal(result, &resp)
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

	repoScanStruct := []RepoScanResult{}

	// Convert bytes to structure
	err = json.Unmarshal(result, &repoScanStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON for %q in result", repo)
	}

	// Convert back to json byte result
	jsonResult, err := json.Marshal(repoScanStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON for %q in result", repo)
	}
	return jsonResult, nil
}

// SaveRepoToDB save repo data.
func (b *Storage) SaveRepoToDB(key string, data []byte) error {
	b.OpenBoltDb()
	defer b.db.Close()

	err := b.db.Update(func(tx *bolt.Tx) error {
		log.Printf("Saving %q to db... \n\n", key)

		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("repo bucket not found")
		}

		// save repo to cache
		err := b.Put([]byte(key), data)
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
