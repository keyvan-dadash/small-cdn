package cache

import (
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

const (
	kJSFile = iota

	kCSSFile

	kSimpleFile
)

const (
	basePath = "/opt"
)

type CacheLog struct {
	gorm.Model
	UserID                uint `gorm:"foreignKey:ID"`
	FileName              string
	FilePath              string
	FileSize              uint64
	FileType              uint8
	DurationOfMinifcation time.Duration
	ConsumedMemory        uint64
}

func CreateCacheLog(fileName string, username string) *CacheLog {
	filePath := filepath.Join(basePath, username, fileName)
	return &CacheLog{
		FileName: fileName,
		FilePath: filePath,
	}
}

func ConvertFileTypeToString(fileType uint8) (error, string) {
	switch fileType {
	case kJSFile:
		{
			return nil, "js"
		}

	case kCSSFile:
		{
			return nil, "css"
		}

	case kSimpleFile:
		{
			return nil, "simple"
		}

	default:
		{
			return ErrUnkownFileType, ""

		}
	}
}
