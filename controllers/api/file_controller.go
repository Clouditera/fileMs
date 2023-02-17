package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"fileMS/model"
	"fileMS/pkg/common"
	"fileMS/pkg/config"
	FS "fileMS/pkg/minio"
	"fileMS/pkg/minio/minio_ext"
	"fileMS/pkg/whitelist"
	"fileMS/services"
	"fmt"
	"github.com/gannicus-w/yunqi_mysql/sqls"
	"github.com/kataras/iris/v12"
	miniov7 "github.com/minio/minio-go/v7"
	"github.com/mlogclub/simple/common/digests"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	gouuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PresignedUploadPartUrlExpireTime = time.Hour * 24 * 7
)

type FileController struct {
	Ctx iris.Context
}

type completedParts []minio_ext.CompletePart

func (a completedParts) Len() int           { return len(a) }
func (a completedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

func (c *FileController) GetChunks() *web.JsonResult {
	var res = -1
	var uuid, uploaded, uploadID, chunks string

	fileMD5 := params.FormValue(c.Ctx, "md5")
	fileName := params.FormValue(c.Ctx, "fileName")
	bucketName := config.Instance.Minio.Bucket

	for {
		// TODO 验证take
		fileChunk := services.FileMsService.Find(sqls.NewCnd().Eq("md5", fileMD5).Eq("file_name", fileName).Desc("id").Limit(1))
		if len(fileChunk) == 0 || fileChunk[0].DeletedAt != nil {
			logrus.Errorf("GetFileChunkByMD5 failed by md5: %s, file name: %s", fileMD5, fileName)
			break
		}

		uuid = fileChunk[0].UUID
		uploaded = strconv.Itoa(fileChunk[0].IsUploaded)
		uploadID = fileChunk[0].UploadID

		prefix := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2]), "/")
		objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], fileChunk[0].FileName), "/")

		isExist, err := isObjectExist(c.Ctx, bucketName, prefix)
		if err != nil {
			logrus.Errorf("isObjectExist failed: %s", err)
			break
		}

		if isExist {
			uploaded = "1"
			if fileChunk[0].IsUploaded != model.FileUploaded {
				logrus.Info("the file has been uploaded but not recorded")

				fileChunk[0].IsUploaded = 1
				if err = services.FileMsService.UpdateColumn(fileChunk[0].Id, "is_uploaded", fileChunk[0].IsUploaded); err != nil {
					logrus.Errorf("UpdateFileChunk failed: %s", err)
				}
			}
			res = 0
			break
		} else {
			uploaded = "0"
			if fileChunk[0].IsUploaded == model.FileUploaded {
				logrus.Info("the file has been recorded but not uploaded")
				fileChunk[0].IsUploaded = 0
				if err = services.FileMsService.UpdateColumn(fileChunk[0].Id, "is_uploaded", fileChunk[0].IsUploaded); err != nil {
					logrus.Errorf("UpdateFileChunk failed: %s", err)
				}
			}
		}

		cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
		if err != nil {
			logrus.Errorf("GetMinioClient failed: %s", err)
			break
		}
		partInfos, err := cl.Client3.ListObjectParts(bucketName, objectName, uploadID)
		if err != nil {
			logrus.Errorf("ListObjectParts failed: %s", err)
			break
		}

		for _, partInfo := range partInfos {
			chunks += strconv.Itoa(partInfo.PartNumber) + "-" + partInfo.ETag + ","
		}

		break
	}

	return web.JsonData(map[string]interface{}{
		"resultCode": strconv.Itoa(res),
		"uuid":       uuid,
		"uploaded":   uploaded,
		"uploadID":   uploadID,
		"chunks":     chunks,
	})
}

// GetNewMultipart md5  fileName  totalChunkCounts  size
func (c *FileController) GetNewMultipart() *web.JsonResult {
	var uuid, uploadID string

	md5 := params.FormValue(c.Ctx, "md5")
	fileName := params.FormValue(c.Ctx, "fileName")

	totalChunkCounts, err := params.FormValueInt(c.Ctx, "totalChunkCounts")
	if nil != err {
		return web.JsonErrorCode(http.StatusBadRequest, "totalChunkCounts is illegal.")
	}

	if totalChunkCounts > minio_ext.MaxPartsCount || totalChunkCounts <= 0 {
		return web.JsonErrorCode(http.StatusBadRequest, "totalChunkCounts is illegal.")
	}

	fileSize, err := params.FormValueInt64(c.Ctx, "size")
	if err != nil {
		return web.JsonErrorCode(http.StatusBadRequest, "size is illegal.")
	}

	if fileSize > minio_ext.MaxMultipartPutObjectSize || fileSize <= 0 {
		return web.JsonErrorCode(http.StatusBadRequest, "size is illegal.")
	}

	fileChunk := services.FileMsService.FindOne(sqls.NewCnd().Eq("md5", md5).Eq("file_name", fileName))
	if fileChunk != nil {
		uuid = fileChunk.UUID
		uploadID, err = newMultiPartUpload(c.Ctx, uuid, fileName)
		if err != nil {
			logrus.Errorf("newMultiPartUpload failed: %s", err)
			return web.JsonErrorCode(http.StatusInternalServerError, "newMultiPartUpload failed.")
		}
		if err = services.FileMsService.Updates(fileChunk.Id, map[string]interface{}{
			"updated_at":      time.Now(),
			"deleted_at":      nil,
			"uuid":            uuid,
			"md5":             md5,
			"is_uploaded":     0,
			"upload_id":       uploadID,
			"total_chunks":    totalChunkCounts,
			"size":            fileSize,
			"file_name":       fileName,
			"completed_parts": "",
		}); err != nil {
			logrus.Errorf("InsertFileChunk failed: %s", err)
			return web.JsonErrorCode(http.StatusInternalServerError, "InsertFileChunk(update) failed.")
		}
	} else {
		uuid = gouuid.NewV4().String()
		uploadID, err = newMultiPartUpload(c.Ctx, uuid, fileName)
		if err != nil {
			logrus.Errorf("newMultiPartUpload failed: %s", err)
			return web.JsonErrorCode(http.StatusInternalServerError, "newMultiPartUpload failed.")
		}
		if err = services.FileMsService.Create(&model.FileChunk{
			UUID:        uuid,
			UploadID:    uploadID,
			Md5:         md5,
			Size:        fileSize,
			FileName:    fileName,
			TotalChunks: totalChunkCounts,
		}); err != nil {
			logrus.Errorf("InsertFileChunk failed: %s", err)
			return web.JsonErrorCode(http.StatusInternalServerError, "InsertFileChunk failed.")
		}
	}

	return web.JsonData(map[string]interface{}{
		"uuid":     uuid,
		"uploadID": uploadID,
	})
}

func (c *FileController) GetMultipartUploadUrl() *web.JsonResult {
	var url string
	uuid := params.FormValue(c.Ctx, "uuid")
	uploadID := params.FormValue(c.Ctx, "uploadID")
	fileName := params.FormValue(c.Ctx, "fileName")

	partNumber, err := params.FormValueInt(c.Ctx, "chunkNumber")
	if err != nil {
		return web.JsonErrorCode(http.StatusBadRequest, "chunkNumber is illegal.")
	}

	size, err := params.FormValueInt64(c.Ctx, "size")
	if err != nil {
		return web.JsonErrorCode(http.StatusBadRequest, "size is illegal.")
	}
	if size > minio_ext.MinPartSize {
		return web.JsonErrorCode(http.StatusBadRequest, "size is illegal.")
	}

	url, err = genMultiPartSignedUrl(uuid, uploadID, fileName, partNumber, size)
	if err != nil {
		logrus.Errorf("genMultiPartSignedUrl failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "genMultiPartSignedUrl failed.")
	}

	return web.JsonData(map[string]interface{}{
		"url": url,
	})
}

func (c *FileController) PostCompleteMultipart() *web.JsonResult {
	uuid := c.Ctx.PostValue("uuid")
	uploadID := c.Ctx.PostValue("uploadID")

	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	uploadInfo, err := completeMultiPartUpload(uuid, uploadID, fileChunk.FileName)
	if err != nil {
		logrus.Errorf("completeMultiPartUpload failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "completeMultiPartUpload failed.")
	}

	fileChunk.IsUploaded = model.FileUploaded

	if err = services.FileMsService.UpdateColumn(fileChunk.Id, "is_uploaded", fileChunk.IsUploaded); err != nil {
		logrus.Errorf("UpdateFileChunk failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "UpdateFileChunk failed.")
	}
	if err = services.FileMsVersionService.Create(&model.FileVersion{
		FileChunkId: fileChunk.Id,
		VersionId:   uploadInfo.VersionID,
	}); err != nil {
		logrus.Errorf("Create FileChunk version failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "Create FileChunk version failed.")
	}

	return web.JsonData(map[string]interface{}{
		"file":    fileChunk.FileName,
		"uuid":    uuid,
		"version": uploadInfo.VersionID,
	})
}

func (c *FileController) PostUpdateMultipart() *web.JsonResult {
	uuid := c.Ctx.PostValue("uuid")
	etag := c.Ctx.PostValue("etag")
	chunkNumber := c.Ctx.PostValue("chunkNumber")

	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	fileChunk.CompletedParts += chunkNumber + "-" + strings.Replace(etag, "\"", "", -1) + ","

	if err := services.FileMsService.UpdateColumn(fileChunk.Id, "completed_parts", fileChunk.CompletedParts); err != nil {
		logrus.Errorf("UpdateFileChunk failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "UpdateFileChunk failed.")
	}

	return web.JsonSuccess()
}

// PostUpload 上传文件(不分片) (maybe discard)
func (c *FileController) PostUpload() *web.JsonResult {
	// TODO 用户身份权限识别
	bucket := c.Ctx.PostValue("bucket")
	prefix := c.Ctx.PostValueTrim("prefix") // 二级文件夹eg sdk， doc
	//bucket := c.Ctx.GetHeader("bucket")

	file, info, err := c.Ctx.FormFile("file") // 32M
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 文件 %s 上传失败", common.ServiceFS+common.ModuleFS+common.ErrBusUpload, info.Filename)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleFS+common.ErrBusUpload, "上传文件失败")
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// TODO 判断bucket是否存在; bucket访问权限等
	if bucket == "" {
		return web.JsonErrorCode(common.ServiceFS+common.ModuleFS+common.ErrPramNull, "参数为空")
	}
	exists, errBucketExists := FS.ClientMinio.Client1.BucketExists(c.Ctx, bucket)
	if !exists {
		logrus.Errorf("errCode: %d, errDescriotion: bucket %s 不存在", common.ServiceFS+common.ModuleFS+common.ErrNotFound, bucket)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrNotFound, "bucket不存在")
	} else if errBucketExists != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 查询存储桶状态异常 %v ", common.ServiceFS+common.ModuleMinio+common.ErrPermission, errBucketExists)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrPermission, "查询存储桶状态异常")
	}

	begin := time.Now()

	n, err := FS.ClientMinio.Client1.PutObject(c.Ctx, bucket, prefix+info.Filename, file, info.Size, miniov7.PutObjectOptions{
		ContentType: "application/octet-stream",
		PartSize:    1024 * 1024 * 16,
	})

	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 文件 %s 上传失败", common.ServiceFS+common.ModuleMinio+common.ErrBusUpload, info.Filename)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrBusUpload, "上传文件失败")
	}

	end := time.Now().Sub(begin)
	fmt.Println("花费时间: ", end)

	logrus.Infof("文件 %s 上传成功", info.Filename)

	return web.JsonData(map[string]interface{}{
		"file": info.Filename,
		"size": n,
	})
}

func (c *FileController) PutUpdateContent() *web.JsonResult {
	uuid := c.Ctx.PostValue("uuid")
	content := c.Ctx.PostValue("content")
	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	if whitelist.CheckList(fileChunk.FileName) {
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrForbidden, "modification of the preset file is prohibited")
	}

	bucketName := config.Instance.Minio.Bucket
	// 重新计算Md5
	reader := bytes.NewBufferString(content)
	md5 := digests.MD5(content)

	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileChunk.FileName), "/")

	uploadInfo, err := FS.ClientMinio.Client1.PutObject(c.Ctx, bucketName, objectName, reader, int64(reader.Len()), miniov7.PutObjectOptions{})
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 更新文件内容失败 %v ", common.ServiceFS+common.ModuleMinio+common.ErrBusUpdateContent, err)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrBusUpdateContent, "更新文件内容失败")
	}

	if err = services.FileMsService.Updates(fileChunk.Id, map[string]interface{}{
		"md5":        md5,
		"size":       uploadInfo.Size,
		"deleted_at": nil,
	}); err != nil {
		logrus.Errorf("InsertFileChunk failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "InsertFileChunk failed.")
	}

	if err = services.FileMsVersionService.Create(&model.FileVersion{
		FileChunkId: fileChunk.Id,
		VersionId:   uploadInfo.VersionID,
	}); err != nil {
		logrus.Errorf("Create FileChunk version failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "Create FileChunk version failed.")
	}

	return web.JsonData(map[string]interface{}{
		"uuid":    uuid,
		"file":    fileChunk.FileName,
		"size":    uploadInfo.Size,
		"md5":     md5,
		"version": uploadInfo.VersionID,
	})
}

// GetContent 获取文本文件内容 (文本不能太大)
func (c *FileController) GetContent() *web.JsonResult {
	uuid := params.FormValue(c.Ctx, "uuid")
	versionId := params.FormValue(c.Ctx, "versionId")

	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	if fileChunk.DeletedAt != nil {
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrNotFound, "file has been deleted.")
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileChunk.FileName), "/")

	var opts miniov7.GetObjectOptions
	if "" != versionId {
		opts = miniov7.GetObjectOptions{
			VersionID: versionId,
		}
	}
	object, err := FS.ClientMinio.Client1.GetObject(c.Ctx, bucketName, objectName, opts)
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 下载文件失败 %v ", common.ServiceFS+common.ModuleMinio+common.ErrBusDownload, err)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrBusDownload, "下载文件失败")
	}
	defer object.Close()

	fileBytes, err := io.ReadAll(object)
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 下载文件失败 %v ", common.ServiceFS+common.ModuleMinio+common.ErrInternal, err)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrInternal, "下载文件失败")
	}
	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	return web.JsonData(
		map[string]interface{}{
			"uuid": uuid,
			"file": fileChunk.FileName,
			"data": encoded,
		})
}

// GetDownload return download url
func (c *FileController) GetDownload() *web.JsonResult {
	uuid := params.FormValue(c.Ctx, "uuid")
	name := params.FormValue(c.Ctx, "name") // 期望下载后修改的文件名
	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	if fileChunk.DeletedAt != nil {
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrNotFound, "file has been deleted.")
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileChunk.FileName), "/")

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	if "" != name {
		//reqParams.Set("response-content-disposition", "attachment; filename=\"11111.txt\"")
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%s", name))
	}

	url, err := FS.ClientMinio.Client1.PresignedGetObject(c.Ctx, bucketName, objectName, time.Second*1000, reqParams)
	if nil != err {
		logrus.Errorf("PresignedGetObject failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "PresignedGetObject failed.")
	}

	return web.JsonData(
		map[string]interface{}{
			"uuid": uuid,
			"url":  url.String(),
			"file": fileChunk.FileName,
		})
}

// Delete 删除文件 可指定version
func (c *FileController) Delete() *web.JsonResult {
	// TODO 用户身份权限识别
	uuid := c.Ctx.PostValue("uuid")

	fileChunk := services.FileMsService.Take("uuid", uuid)
	if fileChunk == nil {
		logrus.Errorf("GetFileChunkByUUID failed by uuid: %s", uuid)
		return web.JsonErrorCode(http.StatusInternalServerError, "GetFileChunkByUUID failed.")
	}

	if whitelist.CheckList(fileChunk.FileName) {
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrForbidden, "modification of the preset file is prohibited")
	}

	_, err := deleteObject(c.Ctx, uuid, fileChunk.FileName)
	if err != nil {
		logrus.Errorf("errCode: %d, errDescriotion: 文件 %s 删除失败", common.ServiceFS+common.ModuleMinio+common.ErrBusDelete, uuid)
		return web.JsonErrorCode(common.ServiceFS+common.ModuleMinio+common.ErrBusDelete, "删除文件失败")
	}

	deleteAt := time.Now()
	if err = services.FileMsService.UpdateColumn(fileChunk.Id, "deleted_at", deleteAt); err != nil {
		logrus.Errorf("UpdateFileChunk failed: %s", err)
		return web.JsonErrorCode(http.StatusInternalServerError, "UpdateFileChunk failed.")
	}
	versionArr := services.FileMsVersionService.Find(sqls.NewCnd().Eq("file_chunk_id", fileChunk.Id).Desc("id").Limit(1))
	if len(versionArr) > 0 && versionArr[0].DeletedAt == nil {
		err = services.FileMsVersionService.UpdateColumn(versionArr[0].Id, "deleted_at", deleteAt)
		if err != nil {
			logrus.Errorf("UpdateFileChunkVersion failed: %s", err)
			return web.JsonErrorCode(http.StatusInternalServerError, "UpdateFileChunkVersion failed.")
		}
	}

	return web.JsonData(
		map[string]interface{}{
			"uuid":      uuid,
			"delete_at": deleteAt,
			"file":      fileChunk.FileName,
		})
}

func (c *FileController) GetVersions() *web.JsonResult {
	uuid := params.FormValue(c.Ctx, "uuid")
	bucketName := config.Instance.Minio.Bucket
	prefix := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid), "/")

	objectCh := FS.ClientMinio.Client1.ListObjects(c.Ctx, bucketName, miniov7.ListObjectsOptions{
		Prefix:       prefix,
		WithVersions: true,
		Recursive:    true,
	})

	var ObjectArr []string
	for object := range objectCh {
		if object.Err == nil {
			ObjectArr = append(ObjectArr, object.VersionID)
		}
	}

	return web.JsonData(ObjectArr)
}

func isObjectExist(ctx context.Context, bucketName string, prefix string) (bool, error) {
	isExist := false
	doneCh := make(chan struct{})
	defer close(doneCh)

	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return isExist, err
	}

	objectCh := cl.Client1.ListObjects(ctx, bucketName, miniov7.ListObjectsOptions{
		Prefix:       prefix,
		WithVersions: true,
		Recursive:    true,
	})
	for object := range objectCh {
		if object.Err != nil {
			logrus.Errorf("ListObjects failed: %s", object.Err)
			return isExist, object.Err
		}
		isExist = true
		break
	}

	return isExist, nil
}

func newMultiPartUpload(ctx iris.Context, uuid, fileName string) (string, error) {
	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return "", err
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileName), "/")

	return cl.Client2.NewMultipartUpload(ctx, bucketName, objectName, miniov7.PutObjectOptions{})
}

func genMultiPartSignedUrl(uuid string, uploadId string, fileName string, partNumber int, partSize int64) (string, error) {
	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return "", err
	}

	bucketName := config.Instance.Minio.Bucket
	//objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, path.Join(uuid[0:1], uuid[1:2], uuid)), "/")
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileName), "/")

	return cl.Client3.GenUploadPartSignedUrl(uploadId, bucketName, objectName, partNumber, partSize, PresignedUploadPartUrlExpireTime, config.Instance.Minio.Location)

}

func completeMultiPartUpload(uuid string, uploadID, fileName string) (minio_ext.UploadInfo, error) {
	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return minio_ext.UploadInfo{}, err
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileName), "/")

	partInfos, err := cl.Client3.ListObjectParts(bucketName, objectName, uploadID)
	if err != nil {
		logrus.Errorf("ListObjectParts failed: %s", err)
		return minio_ext.UploadInfo{}, err
	}

	var complMultipartUpload minio_ext.CompleteMultipartUpload
	for _, partInfo := range partInfos {
		complMultipartUpload.Parts = append(complMultipartUpload.Parts, minio_ext.CompletePart{
			PartNumber: partInfo.PartNumber,
			ETag:       partInfo.ETag,
		})
	}

	// Sort all completed parts.
	sort.Sort(completedParts(complMultipartUpload.Parts))

	return cl.Client3.CompleteMultipartUpload(bucketName, objectName, uploadID, complMultipartUpload.Parts)
}

func deleteObject(ctx iris.Context, uuid, fileName string) (string, error) {
	version := params.FormValue(ctx, "version")
	cl, err := FS.ClientMinio.GetMinioClient(config.Instance)
	if err != nil {
		logrus.Errorf("GetMinioClient failed: %s", err)
		return "", err
	}

	bucketName := config.Instance.Minio.Bucket
	objectName := strings.TrimPrefix(path.Join(config.Instance.Minio.BasePath, uuid[0:1], uuid[1:2], uuid, fileName), "/")

	err = cl.Client1.RemoveObject(ctx, bucketName, objectName, miniov7.RemoveObjectOptions{
		VersionID:        version,
		GovernanceBypass: true,
	})
	if err != nil {
		logrus.Errorf("deleteObject failed: %s", err)
		return "", err
	}

	return uuid, nil
}
