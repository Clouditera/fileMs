package config

import (
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Instance *Config

type Config struct {
	Env      string `yaml:"Env"`      // 环境：prod、dev
	Port     string `yaml:"Port"`     // 端口
	GrpcPort string `yaml:"GrpcPort"` // grpc
	LogFile  string `yaml:"LogFile"`  // 日志文件
	ShowSql  bool   `yaml:"ShowSql"`  // 是否显示日志

	// minio
	Minio struct {
		Url             string `yaml:"Url"`
		AccessKeyID     string `yaml:"AccessKeyID"`
		SecretAccessKey string `yaml:"SecretAccessKey"`
		UseSSL          bool   `yaml:"UseSSL"`
		Bucket          string `yaml:"Bucket"`
		BasePath        string `yaml:"BasePath"`
		Location        string `yaml:"Location"`
	} `yaml:"Minio"`

	// mysql数据库配置
	DB struct {
		Url             string        `yaml:"Url"`
		MaxIdleConns    int           `yaml:"MaxIdleConns"`
		MaxOpenConns    int           `yaml:"MaxOpenConns"`
		MaxLifetimeConn time.Duration `yaml:"MaxLifetimeConn"`
	} `yaml:"DB"`
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		logrus.Error(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		logrus.Error(err)
	}
	return Instance
}
