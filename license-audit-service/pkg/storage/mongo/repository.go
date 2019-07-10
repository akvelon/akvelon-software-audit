package mongo

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage is a mingodb storage realization
type Storage struct {
	db      *mgo.Database
	session *mgo.Session
}

const (
	hosts          = "audit-mongo:27017"
	database       = "auditdb"
	username       = ""
	password       = ""
	repoCollection = "Repository"
	metaCollection = "Meta"

	maxNumOfRecentItems         = 5
	recentMetaCollectionDocName = "Recent"
)

// InitStorage initialize DB for usage.
func (b *Storage) InitStorage() error {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return err
	}

	b.db = s.DB(database)
	b.session = s

	b.session.SetSafe(&mgo.Safe{})

	return nil
}

// GetRecentlyViewed get list of recent items.
func (b *Storage) GetRecentlyViewed() ([]string, error) {
	result := MetaItem{}
	col := b.db.C(metaCollection)
	query := col.Find(bson.M{"name": recentMetaCollectionDocName})

	count, err := query.Count()
	if err != nil {
		log.Println("GetRecentlyViewed error: ", err)
		return nil, err
	}

	if count == 0 {
		return make([]string, 0), nil
	}

	err = query.One(&result)
	if err != nil {
		log.Println("GetRecentlyViewed  error: ", err)
		return nil, err
	}

	// convert to string array and return
	var recentRepos = make([]string, len(result.Recent))
	var j = len(result.Recent) - 1
	for _, r := range result.Recent {
		recentRepos[j] = r.Repo
		j--
	}
	return recentRepos, nil
}

// UpdateRecentlyViewed updated list of recent items.
func (b *Storage) UpdateRecentlyViewed(repo string) error {
	result := MetaItem{Name: recentMetaCollectionDocName}
	col := b.db.C(metaCollection)
	query := col.Find(bson.M{"name": recentMetaCollectionDocName})
	count, err := query.Count()
	if err != nil {
		log.Println("UpdateRecentlyViewed error: ", err)
		return err
	}

	if count > 0 {
		err = query.One(&result)
	} else {
		result.Recent = make([]RecentItem, 0)
	}

	if err != nil {
		log.Println("UpdateRecentlyViewed error: ", err)
		return err
	}

	result.Recent = append(result.Recent, RecentItem{Repo: repo})
	if len(result.Recent) > maxNumOfRecentItems {
		// trim recent if it's grown to over maxNumOfRecentItems
		result.Recent = (result.Recent)[1 : maxNumOfRecentItems+1]
	}

	_, err = col.Upsert(
		bson.M{"name": recentMetaCollectionDocName},
		bson.M{"$set": result},
	)

	if err != nil {
		return err
	}
	return nil
}

// GetRepoFromDB returns repo data if exists.
func (b *Storage) GetRepoFromDB(repo string) ([]RepoScanResult, error) {
	log.Printf("Getting %q data from db... \n\n", repo)
	rsi := RepoScanItem{}
	col := b.db.C(repoCollection)
	err := col.Find(bson.M{"repo": repo}).One(&rsi)
	if err != nil {
		return nil, err
	}
	log.Printf("Scan results: %s", rsi.Results)
	return rsi.Results, nil
}

// SaveRepoToDB save repo data.
func (b *Storage) SaveRepoToDB(repo string, data []RepoScanResult) error {
	log.Printf("Saving %q to db... \n\n", repo)
	col := b.db.C(repoCollection)

	rsi := &RepoScanItem{
		Repo:    repo,
		Results: data,
	}
	_, err := col.Upsert(
		bson.M{"repo": repo},
		bson.M{"$set": rsi},
	)
	if err != nil {
		log.Printf("FAILED to save to db... %s \n\n", err)
		return err
	}
	return nil
}
