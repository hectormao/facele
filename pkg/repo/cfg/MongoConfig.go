package cfg

type MongoConfig interface {
	getURL() string
	getDatabase() string
	getTimeout() int
}
