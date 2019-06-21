package cfg

type MongoConfig interface {
	GetURL() string
	GetDatabase() string
	GetTimeout() int
}
