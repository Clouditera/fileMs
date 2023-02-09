package server

// grpc server
import (
	"context"
	"errors"
	"fileMS/pkg/config"
	FS "fileMS/pkg/minio"
	"fileMS/services"
	"fmt"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type fileMsServer struct{}

func (f fileMsServer) GetFile(ctx context.Context, request *FileRequest) (*FileResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
		return nil, fmt.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(request.Uuid[0:1], request.Uuid[1:2], request.Uuid)), "/")
	url, err := FS.ClientMinio.Client1.PresignedGetObject(bucketName, objectName, time.Second*1000, url.Values{})
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

func (f fileMsServer) DelFile(ctx context.Context, request *FileRequest) (*FileResponse, error) {
	if nil == request || "" == request.Uuid {
		return nil, errors.New("input is null")
	}

	fileChunk := services.FileMsService.Take("uuid", request.Uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
		return nil, fmt.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(request.Uuid[0:1], request.Uuid[1:2], request.Uuid)), "/")
	if err := FS.ClientMinio.Client1.RemoveObject(bucketName, objectName); nil != err {
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
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
		return nil, fmt.Errorf("GetFileChunkByUUID failed by uuid: %s", request.Uuid)
	}
	if nil != fileChunk.DeletedAt {
		return nil, fmt.Errorf("file %s is deleted", request.Uuid)
	}

	files := services.FileMsService.Find(sqls.NewCnd().Where("file_name", fileChunk.FileName))
	if 0 == len(files) {
		return nil, fmt.Errorf("ListFileVersions failed by %s", request.Uuid)
	}
	var versions []string
	for k, v := range files {
		versions = append(versions, "V"+strconv.Itoa(k)+" "+v.CreatedAt.GoString())
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
