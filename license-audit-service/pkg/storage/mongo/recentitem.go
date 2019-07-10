package mongo

// MetaItem represents collection of recent items.
type MetaItem struct {
	Name   string       `bson:"name" json:"name"`
	Recent []RecentItem `bson:"recent" json:"recent"`
}

// RecentItem represents structure for storing recent item related info.
type RecentItem struct {
	Repo string `bson:"repo" json:"repo"`
}