package minio_ext

import (
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v6/pkg/encrypt"
)

// RetentionMode - object retention mode.
type RetentionMode string

func (r RetentionMode) String() string {
	return string(r)
}

// PutObjectOptions represents options specified by user for PutObject call
type PutObjectOptions struct {
	UserMetadata            map[string]string
	Progress                io.Reader
	ContentType             string
	ContentEncoding         string
	ContentDisposition      string
	ContentLanguage         string
	CacheControl            string
	Mode                    *RetentionMode
	RetainUntilDate         *time.Time
	ServerSideEncryption    encrypt.ServerSide
	NumThreads              uint
	StorageClass            string
	WebsiteRedirectLocation string
	PartSize                uint64
}

// getNumThreads - gets the number of threads to be used in the multipart
// put object operation
func (opts PutObjectOptions) getNumThreads() (numThreads int) {
	if opts.NumThreads > 0 {
		numThreads = int(opts.NumThreads)
	} else {
		numThreads = totalWorkers
	}
	return
}

// Header - constructs the headers from metadata entered by user in
// PutObjectOptions struct
func (opts PutObjectOptions) Header() (header http.Header) {
	header = make(http.Header)

	if opts.ContentType != "" {
		header["Content-Type"] = []string{opts.ContentType}
	} else {
		header["Content-Type"] = []string{"application/octet-stream"}
	}
	if opts.ContentEncoding != "" {
		header["Content-Encoding"] = []string{opts.ContentEncoding}
	}
	if opts.ContentDisposition != "" {
		header["Content-Disposition"] = []string{opts.ContentDisposition}
	}
	if opts.ContentLanguage != "" {
		header["Content-Language"] = []string{opts.ContentLanguage}
	}
	if opts.CacheControl != "" {
		header["Cache-Control"] = []string{opts.CacheControl}
	}

	if opts.Mode != nil {
		header["x-amz-object-lock-mode"] = []string{opts.Mode.String()}
	}
	if opts.RetainUntilDate != nil {
		header["x-amz-object-lock-retain-until-date"] = []string{opts.RetainUntilDate.Format(time.RFC3339)}
	}

	if opts.ServerSideEncryption != nil {
		opts.ServerSideEncryption.Marshal(header)
	}
	if opts.StorageClass != "" {
		header[amzStorageClass] = []string{opts.StorageClass}
	}
	if opts.WebsiteRedirectLocation != "" {
		header[amzWebsiteRedirectLocation] = []string{opts.WebsiteRedirectLocation}
	}
	for k, v := range opts.UserMetadata {
		if !isAmzHeader(k) && !isStandardHeader(k) && !isStorageClassHeader(k) {
			header["X-Amz-Meta-"+k] = []string{v}
		} else {
			header[k] = []string{v}
		}
	}
	return
}
