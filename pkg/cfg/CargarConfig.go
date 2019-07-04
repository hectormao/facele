package cfg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type FaceleConfigType struct {
	WebServer         WebServerType     `yaml:"webServer"`
	MongoConfig       MongoConfigType   `yaml:"mongoConfig"`
	WebServiceDian    WebServerDianType `yaml:"webServiceDian"`
	EnvioDianQueue    QueueType         `yaml:"envioDianQueue"`
	NotificacionQueue QueueType         `yaml:"notificacionQueue"`
}

type WebServerDianType struct {
	Url string `yaml:"url"`
}

type WebServerType struct {
	Port   int    `yaml:"port"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type MongoConfigType struct {
	URL      string `yaml:"url"`
	Timeout  int    `yaml:"timeout"`
	Database string `yaml:"database"`
}

func (cfg MongoConfigType) GetURL() string {
	return cfg.URL
}

func (cfg MongoConfigType) GetDatabase() string {
	return cfg.Database
}

func (cfg MongoConfigType) GetTimeout() int {
	return cfg.Timeout
}

type QueueType struct {
	Name string `yaml:"name"`
}

func CargarConfig(configPath string) (FaceleConfigType, error) {

	var config FaceleConfigType

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil

}
