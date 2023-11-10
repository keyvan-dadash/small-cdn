package cache

import (
	"bytes"
	"io"

	"github.com/tdewolff/minify/v2"
)

type JSCacheFile struct {
	minifier *minify.M
	reader   io.Reader
	buffer   *bytes.Buffer
}

func CreateJSCacheFile(minifier *minify.M) *JSCacheFile {
	return &JSCacheFile{
		minifier: minifier,
		buffer:   new(bytes.Buffer),
	}
}

func (j *JSCacheFile) Write(writer io.Writer) error {
	_, err := writer.Write(j.buffer.Bytes())
	return err
}

func (j *JSCacheFile) Read(reader io.Reader) error {
	j.reader = reader
	return nil
}

func (j *JSCacheFile) Process() error {
	if err := j.minifier.Minify("application/javascript", j.buffer, j.reader); err != nil {
		return err
	}

	return nil
}

func (j *JSCacheFile) GetCacheFileType() uint8 {
	return kJSFile
}
