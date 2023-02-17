/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2017 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package minio_ext

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/minio/minio-go/v6/pkg/s3utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

const amzVersionID = "X-Amz-Version-Id"

// CompleteMultipartUpload - Concatenate uploaded parts and commit to an object.
func (c *Client) CompleteMultipartUpload(bucket, object, uploadID string, parts []CompletePart) (UploadInfo, error) {
	res, err := c.completeMultipartUpload(context.Background(), bucket, object, uploadID, CompleteMultipartUpload{
		Parts: parts,
	}, PutObjectOptions{})
	return res, err
}

// completeMultipartUpload - Completes a multipart upload by assembling previously uploaded parts.
func (c *Client) completeMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string,
	complete CompleteMultipartUpload, opts PutObjectOptions,
) (UploadInfo, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return UploadInfo{}, err
	}
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return UploadInfo{}, err
	}

	// Initialize url queries.
	urlValues := make(url.Values)
	urlValues.Set("uploadId", uploadID)
	// Marshal complete multipart body.
	completeMultipartUploadBytes, err := xml.Marshal(complete)
	if err != nil {
		return UploadInfo{}, err
	}

	// Instantiate all the complete multipart buffer.
	completeMultipartUploadBuffer := bytes.NewReader(completeMultipartUploadBytes)
	reqMetadata := requestMetadata{
		bucketName:       bucketName,
		objectName:       objectName,
		queryValues:      urlValues,
		contentBody:      completeMultipartUploadBuffer,
		contentLength:    int64(len(completeMultipartUploadBytes)),
		contentSHA256Hex: sum256Hex(completeMultipartUploadBytes),
		customHeader:     opts.Header(),
	}

	// Execute POST to complete multipart upload for an objectName.
	resp, err := c.executeMethod(ctx, http.MethodPost, reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return UploadInfo{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return UploadInfo{}, httpRespToErrorResponse(resp, bucketName, objectName)
		}
	}

	// Read resp.Body into a []bytes to parse for Error response inside the body
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return UploadInfo{}, err
	}
	// Decode completed multipart upload response on success.
	completeMultipartUploadResult := completeMultipartUploadResult{}
	err = xmlDecoder(bytes.NewReader(b), &completeMultipartUploadResult)
	if err != nil {
		// xml parsing failure due to presence an ill-formed xml fragment
		return UploadInfo{}, err
	} else if completeMultipartUploadResult.Bucket == "" {
		// xml's Decode method ignores well-formed xml that don't apply to the type of value supplied.
		// In this case, it would leave completeMultipartUploadResult with the corresponding zero-values
		// of the members.

		// Decode completed multipart upload response on failure
		completeMultipartUploadErr := ErrorResponse{}
		err = xmlDecoder(bytes.NewReader(b), &completeMultipartUploadErr)
		if err != nil {
			// xml parsing failure due to presence an ill-formed xml fragment
			return UploadInfo{}, err
		}
		return UploadInfo{}, completeMultipartUploadErr
	}

	return UploadInfo{
		Bucket:    completeMultipartUploadResult.Bucket,
		Key:       completeMultipartUploadResult.Key,
		ETag:      trimEtag(completeMultipartUploadResult.ETag),
		VersionID: resp.Header.Get(amzVersionID),
		Location:  completeMultipartUploadResult.Location,

		ChecksumSHA256: completeMultipartUploadResult.ChecksumSHA256,
		ChecksumSHA1:   completeMultipartUploadResult.ChecksumSHA1,
		ChecksumCRC32:  completeMultipartUploadResult.ChecksumCRC32,
		ChecksumCRC32C: completeMultipartUploadResult.ChecksumCRC32C,
	}, nil
}
