package configext

type MongoConfig struct {
	Url        string `json:"url"`
	Db         string `json:"db"`
	Collection string `json:"collection"`
}
