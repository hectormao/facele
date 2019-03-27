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
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
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
