package minio

import (
	"github.com/minio/minio-go/v6"
	"log"
	"testing"
)

// 使用说明文档： https://clouditera.feishu.cn/docx/GFOpdfF5qojSgzxwrAmcd4venlG

func TestMinIO(t *testing.T) {
	endpoint := "192.168.31.18:9000"
	accessKeyID := "minio"
	secretAccessKey := "minio123"
	useSSL := false
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)
	}
	log.Printf("创建 MinIO 客户端成功")

	// create a new bucket
	bucketName := "mymusic"
	//bucketName := "clouditera"
	location := "cn-north-1"
	err = minioClient.MakeBucket(bucketName, location) //Region: location
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("存储桶 %s 已经存在", bucketName)
		} else {
			log.Fatalln("查询存储桶状态异常", errBucketExists)
		}
	}
	log.Printf("创建存储桶 %s 成功", bucketName)

	// Upload file
	objectName := "README.md"
	filePath := "D:\\Code\\go\\tmp\\v1.5\\fileMSPlatform\\README.md"
	contentType := "application/zip"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("上传文件失败", err)
	}

	log.Printf("上传文件 %s成功, ----> %v\n", objectName, n)

	// minioClient.FGetObject()

	//file, info, err := c.Ctx.FormFile("license")
	//if err != nil {
	//
	//}
	//defer func(file multipart.File) {
	//	err := file.Close()
	//	if err != nil {
	//
	//	}
	//}(file)
	//
	//minioClient.PutObject()

}
