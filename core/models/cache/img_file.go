package cache

import "io"

type IMGCacheFile struct {
}

func CreateIMGCacheFile() *IMGCacheFile {
	return &IMGCacheFile{}
}

func (j *IMGCacheFile) Write(writer io.Writer) error {
	return nil
}

func (j *IMGCacheFile) Read(reader io.Reader) error {
	return nil
}

func (j *IMGCacheFile) Process() error {
	return nil
}

func (j *IMGCacheFile) GetCacheFileType() uint8 {
	return kJSFile
}
