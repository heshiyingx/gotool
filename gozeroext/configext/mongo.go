package configext

type MongoConfig struct {
	Url        string `json:"url,optional"`
	Db         string `json:"db,optional"`
	Collection string `json:"collection,optional"`
}
