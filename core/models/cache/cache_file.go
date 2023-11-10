package cache

import (
	"errors"
	"io"
)

var (
	ErrUnkownFileType = errors.New("Provided CacheFile is unkown")
)

type CacheFile interface {
	Read(io.Reader) error
	Write(io.Writer) error
	Process() error
	GetCacheFileType() uint8
}

func CreateCacheFileFactory(cacheLog *CacheLog, fileTypeStr string) (error, CacheFile) {
	var cacheFile CacheFile

	switch fileTypeStr {
	case "js":
		{
			cacheLog.FileType = kJSFile
			cacheFile = CreateJSCacheFile()
		}

	case "css":
		{
			cacheLog.FileType = kCSSFile
			cacheFile = CreateCSSCacheFile()
		}

	case "image":
		{
			cacheLog.FileType = kImgFile
			cacheFile = CreateIMGCacheFile()
		}
	default:
		{
			return ErrUnkownFileType, nil
		}
	}
	return nil, cacheFile
}
