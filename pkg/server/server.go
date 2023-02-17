package server

// grpc server
import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fileMS/model"
	"fileMS/pkg/common"
	"fileMS/pkg/config"
	FS "fileMS/pkg/minio"
	"fileMS/pkg/whitelist"
	"fileMS/services"
	"fmt"
	"github.com/gannicus-w/yunqi_mysql/common/digests"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	miniov7 "github.com/minio/minio-go/v7"
	gouuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/url"
	"path"
	"strings"
	"time"
)

type fileMsServer struct{}

func (f fileMsServer) CreateFile(ctx context.Context, request *CreateFileRequest) (*FileRequest, error) {
	if nil == request {
		return nil, errors.New("input is null")
	}
	// 计算md5
	reader := bytes.NewBufferString(request.Content)
	md5 := digests.MD5(request.Content)
	uuid := gouuid.NewV4().String()
	uploadID, err := newMultiPartUpload(ctx, uuid, request.Filename)
	if err != nil {
		logrus.Errorf("newMultiPartUpload failed: %s", err)
		return nil, err
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(uuid[0:1], uuid[1:2], uuid, request.Filename)), "/")

	uploadInfo, err := FS.ClientMinio.Client1.PutObject(ctx, bucketName, objectName, reader, int64(reader.Len()), miniov7.PutObjectOptions{})
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 更新文件内容失败 %v ", common.ServiceFS+common.ModuleMinio+common.ErrBusUpdateContent, err)
		return nil, err
	}

	fileChunk := &model.FileChunk{
		UUID:        uuid,
		UploadID:    uploadID,
		Md5:         md5,
		Size:        uploadInfo.Size,
		FileName:    request.Filename,
		TotalChunks: 1,
	}
	if err = services.FileMsService.Create(fileChunk); err != nil {
		logrus.Errorf("InsertFileChunk failed: %s", err)
		return nil, err
	}

	if err = services.FileMsVersionService.Create(&model.FileVersion{
		FileChunkId: fileChunk.Id,
		VersionId:   uploadInfo.VersionID,
	}); err != nil {
		logrus.Errorf("Create FileChunk version failed: %s", err)
		return nil, err
	}

	return &FileRequest{
		Uuid:    uuid,
		Version: uploadInfo.VersionID,
	}, nil
}

func (f fileMsServer) UpdateFileContent(ctx context.Context, request *UpdateFileRequest) (*FileRequest, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}
	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFile failed by uuid: %s, when qurey in mysql", request.Uuid)
		return nil, fmt.Errorf("GetFile failed by uuid: %s, when qurey in mysql", request.Uuid)
	}

	if whitelist.CheckList(fileChunk.FileName) {
		return nil, fmt.Errorf("modification of the preset file is prohibited")
	}

	// 重新计算md5
	reader := bytes.NewBufferString(request.Content)
	md5 := digests.MD5(request.Content)
	if 0 == strings.Compare(md5, fileChunk.Md5) {
		return &FileRequest{
			Uuid: fileChunk.UUID,
		}, nil
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(fileChunk.UUID[0:1], fileChunk.UUID[1:2], fileChunk.UUID, fileChunk.FileName)), "/")

	uploadInfo, err := FS.ClientMinio.Client1.PutObject(ctx, bucketName, objectName, reader, int64(reader.Len()), miniov7.PutObjectOptions{})
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 更新文件内容失败 %v ", common.ServiceFS+common.ModuleMinio+common.ErrBusUpdateContent, err)
		return nil, err
	}

	if err = services.FileMsService.Updates(fileChunk.Id, map[string]interface{}{
		"md5":  md5,
		"size": 5,
	}); err != nil {
		logrus.Errorf("InsertFileChunk failed: %s", err)
		return nil, err
	}

	if err = services.FileMsVersionService.Create(&model.FileVersion{
		FileChunkId: fileChunk.Id,
		VersionId:   uploadInfo.VersionID,
	}); err != nil {
		logrus.Errorf("Create FileChunk version failed: %s", err)
		return nil, err
	}

	return &FileRequest{
		Uuid:    fileChunk.UUID,
		Version: uploadInfo.VersionID,
	}, nil
}

func (f fileMsServer) GetFile(ctx context.Context, request *FileRequest) (*FileResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFile failed by uuid: %s, when qurey in mysql", request.Uuid)
		return nil, fmt.Errorf("GetFile failed by uuid: %s, when qurey in mysql", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(request.Uuid[0:1], request.Uuid[1:2], request.Uuid, fileChunk.FileName)), "/")

	url, err := FS.ClientMinio.Client1.PresignedGetObject(ctx, bucketName, objectName, time.Second*1000, url.Values{})
	if nil != err {
		logrus.Errorf("PresignedGetObject failed by uuid: %s", request.Uuid)
		return nil, fmt.Errorf("PresignedGetObject failed by uuid: %s", request.Uuid)
	}

	return &FileResponse{
		Uuid: fileChunk.UUID,
		File: fileChunk.FileName,
		Url:  url.String(),
	}, nil
}

func (f fileMsServer) GetFileContent(ctx context.Context, request *FileRequest) (*FileContentResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileContent failed by uuid: %s, when qurey in mysql", request.Uuid)
		return nil, fmt.Errorf("GetFileContent failed by uuid: %s, when qurey in mysql", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(request.Uuid[0:1], request.Uuid[1:2], request.Uuid, fileChunk.FileName)), "/")

	object, err := FS.ClientMinio.Client1.GetObject(ctx, bucketName, objectName, miniov7.GetObjectOptions{})
	if err != nil {
		logrus.Errorf("GetFileContent failed by uuid: %s", request.Uuid)
		return nil, fmt.Errorf("GetFileContent failed by uuid: %s", request.Uuid)
	}
	defer object.Close()

	objectInfo, err := object.Stat()
	if err != nil {
		logrus.Errorf("GetFileContent failed by uuid: %s, when object state", request.Uuid)
		return nil, fmt.Errorf("GetFileContent failed by uuid: %s, when object state", request.Uuid)
	}
	// 200M
	if objectInfo.Size >= 1024*1024*200 {
		logrus.Errorf("GetFileContent failed by uuid: %s, the content more than 200M", request.Uuid)
		return nil, fmt.Errorf("GetFileContent failed by uuid: %s, the content more than 200M", request.Uuid)
	}

	fileBytes, err := io.ReadAll(object)
	if err != nil {
		logrus.Errorf("GetFileContent failed by uuid: %s, when read content", request.Uuid)
		return nil, fmt.Errorf("GetFileContent failed by uuid: %s, when read content", request.Uuid)
	}
	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	return &FileContentResponse{
		Uuid:    fileChunk.UUID,
		File:    fileChunk.FileName,
		Content: encoded,
	}, nil
}

func (f fileMsServer) DelFile(ctx context.Context, request *FileRequest) (*FileResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("DelFile failed by uuid: %s, when qurey in mysql", request.Uuid)
		return nil, fmt.Errorf("DelFile failed by uuid: %s, when qurey in mysql", request.Uuid)
	}

	if whitelist.CheckList(fileChunk.FileName) {
		return nil, fmt.Errorf("modification of the preset file is prohibited")
	}

	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(request.Uuid[0:1], request.Uuid[1:2], request.Uuid, fileChunk.FileName)), "/")

	if err := FS.ClientMinio.Client1.RemoveObject(ctx, bucketName, objectName, miniov7.RemoveObjectOptions{}); nil != err {
		logrus.Errorf("deleteObject failed: %s", err)
		return nil, err
	}

	deleteAt := time.Now()
	if err := services.FileMsService.UpdateColumn(fileChunk.Id, "deleted_at", deleteAt); err != nil {
		logrus.Errorf("UpdateFileChunk failed: %s", err)
		return nil, fmt.Errorf("UpdateFileChunk failed: %s", err)
	}

	return &FileResponse{
		Uuid: fileChunk.UUID,
		File: fileChunk.FileName,
	}, nil

}

func (f fileMsServer) ListFileVersions(ctx context.Context, request *FileRequest) (*FileVersionResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("ListFileVersions failed by uuid: %s, when qurey in mysql", request.Uuid)
		return nil, fmt.Errorf("ListFileVersions failed by uuid: %s, when qurey in mysql", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	rets := services.FileMsVersionService.Find(sqls.NewCnd().Eq("file_chunk_id", fileChunk.Id))

	var versions []string
	for _, v := range rets {
		if v.DeletedAt == nil {
			versions = append(versions, v.VersionId)
		}
	}

	return &FileVersionResponse{
		Uuid:    fileChunk.UUID,
		File:    fileChunk.FileName,
		Version: versions,
	}, nil
}

func InitServer(port string) error {
	rpcServer := grpc.NewServer()
	RegisterFileMsServiceServer(rpcServer, new(fileMsServer))

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Errorf("failed set grpc listening port: %s", err)
		return err
	}

	err = rpcServer.Serve(lis)
	if err != nil {
		logrus.Errorf("failed to start grpc server: %s", err)
		return err
	}
	return nil
}

func newMultiPartUpload(ctx context.Context, uuid, fileName string) (string, error) {
	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return "", err
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileName), "/")

	return cl.Client2.NewMultipartUpload(ctx, bucketName, objectName, miniov7.PutObjectOptions{})
}
