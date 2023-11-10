package cache

import "io"

type JSCacheFile struct {
}

func CreateJSCacheFile() *JSCacheFile {
	return &JSCacheFile{}
}

func (j *JSCacheFile) Write(writer io.Writer) error {
	return nil
}

func (j *JSCacheFile) Read(reader io.Reader) error {
	return nil
}

func (j *JSCacheFile) Process() error {
	return nil
}

func (j *JSCacheFile) GetCacheFileType() uint8 {
	return kJSFile
}
