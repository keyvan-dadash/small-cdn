package cache

import (
	"errors"
	"io"
	"regexp"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	ErrUnkownFileType = errors.New("Provided CacheFile is unkown")
)

var (
	minifyOnce sync.Once
	minifier   *minify.M
)

type CacheFile interface {
	Read(io.Reader) error
	Write(io.Writer) error
	Process() error
	GetCacheFileType() uint8
}

func createMinifier() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

func CreateCacheFileFactory(cacheLog *CacheLog, fileTypeStr string) (error, CacheFile) {
	minifyOnce.Do(createMinifier)

	var cacheFile CacheFile

	switch fileTypeStr {
	case "js":
		{
			cacheLog.FileType = kJSFile
			cacheFile = CreateJSCacheFile(minifier)
		}

	case "css":
		{
			cacheLog.FileType = kCSSFile
			cacheFile = CreateCSSCacheFile(minifier)
		}

	case "simple":
		{
			cacheLog.FileType = kSimpleFile
			cacheFile = CreateSimpleCacheFile()
		}

	default:
		{
			return ErrUnkownFileType, nil
		}
	}
	return nil, cacheFile
}
