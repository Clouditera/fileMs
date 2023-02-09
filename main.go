package main

import (
	"fileMS/controllers"
	"fileMS/model"
	"fileMS/pkg/common"
	"fileMS/pkg/config"
	"fileMS/pkg/minio"
	"flag"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configFile = flag.String("config", "./filems.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)

	// gorm配置
	gormConf := &gorm.Config{}

	// 初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		if conf.ShowSql {
			gormConf.Logger = logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
		}
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 启动minio
	if err := minio.InitMinio(conf); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("success to connect minio")
	}

	// 启动mysql
	if err := sqls.Open(conf.DB.Url, gormConf, conf.DB.MaxIdleConns, conf.DB.MaxOpenConns, conf.DB.MaxLifetimeConn*time.Hour, model.Models...); err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("success to connect MySQL [%s]", conf.DB.Url)
	}

	// 启动grpc server
	//if err := grpc.InitServer(conf.GrpcPort); err != nil {
	//	logrus.Error(err)
	//} else {
	//	logrus.Infof("success to start grpc server")
	//}
}

func main() {
	if common.IsProd() {
		// TODO: 预留
	}
	controllers.Router()
}
