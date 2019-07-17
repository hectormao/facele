package cfg

import (
	"io/ioutil"

	repo "github.com/hectormao/facele/pkg/repo/cfg"
	trns "github.com/hectormao/facele/pkg/trns/cfg"
	"gopkg.in/yaml.v2"
)

type FaceleConfigType struct {
	WebServer      WebServerType     `yaml:"webserver"`
	MongoConfig    MongoConfigType   `yaml:"mongo_config"`
	WebServiceDian WebServerDianType `yaml:"webservice_dian"`
	RabbitConfig   RabbitConfigType  `yaml:"rabbit_config"`
}

type RabbitConfigType struct {
	Url               string    `yaml:"url"`
	EnvioDianQueue    QueueType `yaml:"envio_dian_queue"`
	NotificacionQueue QueueType `yaml:"notificacion_queue"`
}

func (cfg RabbitConfigType) GetUrl() string {
	return cfg.Url
}

func (cfg RabbitConfigType) GetEnvioDianQueue() repo.QueueConfig {
	return cfg.EnvioDianQueue
}

func (cfg RabbitConfigType) GetNotificacionQueue() repo.QueueConfig {
	return cfg.NotificacionQueue
}

type WebServerDianType struct {
	Url       string       `yaml:"url"`
	QRCodeUrl string       `yaml:"qrcode_url"`
	SignCert  SignCertType `yaml:"sign_cert"`
}

func (cfg WebServerDianType) GetQRCodeUrl() string {
	return cfg.QRCodeUrl
}

func (cfg WebServerDianType) GetSignCert() trns.SignCertConfig {
	return cfg.SignCert
}

type SignCertType struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func (cfg SignCertType) GetPath() string {
	return cfg.Path
}

func (cfg SignCertType) GetPassword() string {
	return cfg.Password
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

func (cfg QueueType) GetName() string {
	return cfg.Name
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
