package cache

import "io"

type CSSCacheFile struct {
}

func CreateCSSCacheFile() *CSSCacheFile {
	return &CSSCacheFile{}
}

func (j *CSSCacheFile) Write(writer io.Writer) error {
	return nil
}

func (j *CSSCacheFile) Read(reader io.Reader) error {
	return nil
}

func (j *CSSCacheFile) Process() error {
	return nil
}

func (j *CSSCacheFile) GetCacheFileType() uint8 {
	return kJSFile
}
