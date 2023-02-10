package model

import "time"

var Models = []interface{}{
	&FileChunk{},
}

const (
	FileNotUploaded int = iota
	FileUploaded
)

type Model struct {
	Id        int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FileChunk struct {
	//ID			  int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Model

	UUID           string `gorm:"UNIQUE"`
	Md5            string // `gorm:"UNIQUE"`
	IsUploaded     int    `gorm:"DEFAULT 0"` // not uploaded: 0, uploaded: 1
	UploadID       string `gorm:"UNIQUE"`    //minio upload id
	TotalChunks    int
	Size           int64
	FileName       string
	CompletedParts string `gorm:"type:text"` // chunkNumber+etag eg: ,1-asqwewqe21312312.2-123hjkas
}
