package minio

import (
	"context"
	"fileMS/pkg/config"
	"fileMS/pkg/minio/minio_ext"
	miniov7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"sync"
)

// Minio 3个client，本质只有一个，写三个的原本因是有些接口不能直接调用，需要扩展
type Minio struct {
	Client1 *miniov7.Client
	Client2 *miniov7.Core
	Client3 *minio_ext.Client
}

var ClientMinio *Minio

var mutex *sync.Mutex

func init() {
	mutex = new(sync.Mutex)
}

func InitMinio(conf *config.Config) error {
	var err error
	var client1 *miniov7.Client
	var client2 *miniov7.Core
	var client3 *minio_ext.Client

	// Initialize minio client object.
	client1, err = miniov7.New(conf.Minio.Url, &miniov7.Options{
		Creds:  credentials.NewStaticV4(conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, ""),
		Secure: conf.Minio.UseSSL,
	})
	if err != nil {
		logrus.Errorf("Minio client1 初始化失败：%s", err)
	}
	client2, err = miniov7.NewCore(conf.Minio.Url, &miniov7.Options{
		Creds:  credentials.NewStaticV4(conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, ""),
		Secure: conf.Minio.UseSSL,
	})
	if err != nil {
		logrus.Errorf("Minio client2 初始化失败：%s", err)
	}

	client3, err = minio_ext.New(conf.Minio.Url, conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, conf.Minio.UseSSL)
	if err != nil {
		logrus.Errorf("Minio client3 初始化失败：%s", err)
	}

	if err != nil {
		logrus.Errorf("Minio client3 初始化失败：%s", err)
	}

	ClientMinio = &Minio{
		Client1: client1,
		Client2: client2,
		Client3: client3,
	}

	logrus.Info("Minio 初始化成功!")

	return nil
}

func (c *Minio) GetMinioClient(conf *config.Config) (*Minio, error) {
	var err error
	mutex.Lock()

	if nil != c.Client1 && nil != c.Client2 && nil != c.Client3 {
		mutex.Unlock()
		return c, nil
	}

	if nil == c.Client1 {
		c.Client1, err = miniov7.New(conf.Minio.Url, &miniov7.Options{
			Creds:  credentials.NewStaticV4(conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, ""),
			Secure: conf.Minio.UseSSL,
		})
		if nil != err {
			mutex.Unlock()
			logrus.Error(err)
			return nil, err
		}
	}

	if nil == c.Client2 {
		c.Client2, err = miniov7.NewCore(conf.Minio.Url, &miniov7.Options{
			Creds:  credentials.NewStaticV4(conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, ""),
			Secure: conf.Minio.UseSSL,
		})
		if nil != err {
			mutex.Unlock()
			logrus.Error(err)
			return nil, err
		}
	}

	if nil == c.Client3 {
		c.Client3, err = minio_ext.New(conf.Minio.Url, conf.Minio.AccessKeyID, conf.Minio.SecretAccessKey, conf.Minio.UseSSL)
		if nil != err {
			mutex.Unlock()
			logrus.Error(err)
			return nil, err
		}
	}

	mutex.Unlock()

	return c, nil
}

func (c *Minio) CreatBucket(bucketName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Client1.MakeBucket(ctx, bucketName, miniov7.MakeBucketOptions{
		Region:        config.Instance.Minio.Location,
		ObjectLocking: true,
	})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := c.Client1.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			logrus.Infof("存储桶 %s 已经存在", bucketName)
		} else {
			logrus.Errorf("查询存储桶状态异常: %v", errBucketExists)
			return errBucketExists
		}
	} else {
		logrus.Infof("创建存储桶 %s 成功", bucketName)
	}

	return nil
}
